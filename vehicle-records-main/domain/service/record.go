package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"git.abanppc.com/farin-project/vehicle-records/domain/entity"
	"git.abanppc.com/farin-project/vehicle-records/domain/opentelemetry"
	"git.abanppc.com/farin-project/vehicle-records/domain/repository"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/godotenv"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/rabbit"
	"github.com/google/uuid"
	"github.com/mahdimehrabi/uploader"
	"github.com/mahdimehrabi/uploader/minio"
	minio2 "github.com/minio/minio-go/v7"
	"github.com/rabbitmq/amqp091-go"
	ptime "github.com/yaa110/go-persian-calendar"
	"go.opentelemetry.io/otel/metric"
)

var ErrFailedUpdate = errors.New("failed to update record")

type VehicleRecordService struct {
	logger           *slog.Logger
	recordRepo       *repository.VehicleRecordRepository
	photoRepo        *repository.CitizenVehiclePhotoRepository
	fr               uploader.FileRepository
	env              *godotenv.Env
	tehranSiteRepo   *repository.TehranSiteRecordRepository
	rbt              *rabbit.Rabbit
	osDuration       metric.Int64Histogram
	shardariDuration metric.Int64Histogram
	processDuration  metric.Int64Histogram
	minio            *minio.Minio
}

func NewVehicleRecordService(logger *slog.Logger, ringRepo *repository.VehicleRecordRepository, fr uploader.FileRepository,
	env *godotenv.Env, photoRepo *repository.CitizenVehiclePhotoRepository, tehranSiteRepo *repository.TehranSiteRecordRepository,
	rbt *rabbit.Rabbit, telemetry *opentelemetry.OpenTelemetry, minio *minio.Minio) *VehicleRecordService {
	meter := telemetry.Meter.Meter("farin,vehicle_record.RabbitMQVehicleRecordHandler")
	osDuration, err := meter.Int64Histogram("vehicle_record.objectstorage.duration",
		metric.WithDescription("time of vehicle record object storage upload"))
	if err != nil {
		panic(err)
	}
	shardariDuration, err := meter.Int64Histogram("vehicle_record.shardari.duration",
		metric.WithDescription("time of sending to shardari"))
	if err != nil {
		panic(err)
	}
	processDuration, err := meter.Int64Histogram("vehicle_record.store.duration",
		metric.WithDescription("time of process"))
	if err != nil {
		panic(err)
	}

	return &VehicleRecordService{
		logger:           logger.With("layer", "VehicleRecordService"),
		recordRepo:       ringRepo,
		fr:               fr,
		env:              env,
		photoRepo:        photoRepo,
		tehranSiteRepo:   tehranSiteRepo,
		rbt:              rbt,
		osDuration:       osDuration,
		shardariDuration: shardariDuration,
		processDuration:  processDuration,
		minio:            minio,
	}
}

func (s *VehicleRecordService) CreateRecord(ctx context.Context, record *entity.VehicleRecord) (*entity.VehicleRecord, error) {
	startProcessDuration := time.Now().Unix()
	tStart := time.Now()
	lg := s.logger.With("method", "CreateRecord")
	defer func() {
		lg.Info("create record", "duration", time.Since(tStart))
	}()
	vfs := make([]*entity.CitizenVehiclePhoto, len(record.VehiclePhotos))
	bts := make([]*entity.CitizenVehiclePhoto, len(record.VehiclePhotos))
	for i, vf := range record.VehiclePhotos {
		vf.ID = uuid.NewString()
		var record1, record2 entity.CitizenVehiclePhoto = *vf, *vf
		vfs[i] = &record1
		bts[i] = &record2
	}
	record.VehiclePhotos = nil
	record.Sent = false

	// Keep track of uploaded files for cleanup
	uploadedFiles := make([]string, 0)
	var createdPhotoIDs []string // Track created photo IDs for cleanup

	record.Sent = true
	var err error
	record.CitizenPlateNumberNumeric, err = plateReplaceAlphabets(record.CitizenPlateNumber)
	if err != nil {
		lg.Error("failed to convert plate string number to number", "error", err,
			"citizen plate number", record.CitizenPlateNumber, "record id ", record.RecordID)
		return nil, err
	}

	// Get a new instance of ptime.Time using time.Time
	pt := ptime.Now()
	record.ShamsiTime = pt.String()
	createdRecord, err := s.recordRepo.Create(ctx, record)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateRecord) {
			lg.Warn("record already exists", "record_id", record.RecordID)
			return nil, err
		}
		lg.Error("failed to create record", "error", err.Error())
		return nil, err
	}

	// Define cleanup function that handles all rollback operations
	cleanup := func() {
		// Delete the record
		if err := s.recordRepo.Delete(ctx, record.RecordID); err != nil {
			lg.Error("failed to delete record during rollback", "error", err.Error())
		}

		// Bulk delete all created photos
		if len(createdPhotoIDs) > 0 {
			if err := s.photoRepo.BulkDelete(ctx, createdPhotoIDs); err != nil {
				lg.Error("failed to delete photos during rollback", "error", err.Error())
			}
		}

		// Delete all uploaded files
		for _, fname := range uploadedFiles {
			if err := s.fr.DeleteFile(ctx, s.env.MinioVehicleRecordsBucket, fname); err != nil {
				lg.Error("failed to delete file during rollback", "filename", fname, "error", err.Error())
			}
		}
	}

	startOSUploadTime := time.Now().UnixMicro()
	for i, photo := range vfs {
		// Upload main photo
		mainFname := fmt.Sprintf("record_id_%s_photo_id_%s", record.RecordID, photo.ID)

		// Decode the main photo
		decodedLen := base64.StdEncoding.DecodedLen(len(photo.CitizenVehiclePhoto))
		b := make([]byte, decodedLen)
		n, err := base64.StdEncoding.Decode(b, []byte(photo.CitizenVehiclePhoto))
		if err != nil {
			lg.Error("failed to decode citizen photo", "error", err.Error())
			return createdRecord, err
		}
		// Use only the actual decoded bytes
		b = b[:n]

		rs := bytes.NewReader(b)
		err = s.fr.UploadPublicFile(ctx, s.env.MinioVehicleRecordsBucket,
			mainFname,
			"application/octet-stream", rs)
		if err != nil {
			cleanup()
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		uploadedFiles = append(uploadedFiles, mainFname)
		vfs[i].CitizenVehiclePhoto = fmt.Sprintf("%s/%s", s.env.MinioVehicleRecordsBucket, mainFname)

		// Decode the crop photo
		decodedLen = base64.StdEncoding.DecodedLen(len(photo.CitizenVehiclePlateCropPhoto))
		c := make([]byte, decodedLen)
		n, err = base64.StdEncoding.Decode(c, []byte(photo.CitizenVehiclePlateCropPhoto))
		if err != nil {
			lg.Error("failed to decode citizen photo", "error", err.Error())
			return createdRecord, err
		}
		// Use only the actual decoded bytes
		c = c[:n]

		// Upload crop photo
		cropFname := fmt.Sprintf("record_id_%s_crop_photo_id_%s", record.RecordID, photo.ID)
		rs = bytes.NewReader(c)
		err = s.fr.UploadPublicFile(ctx, s.env.MinioVehicleRecordsBucket,
			cropFname,
			"application/octet-stream", rs)
		if err != nil {
			cleanup()
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		uploadedFiles = append(uploadedFiles, cropFname)
		vfs[i].CitizenVehiclePlateCropPhoto = fmt.Sprintf("%s/%s", s.env.MinioVehicleRecordsBucket, cropFname)
	}
	lg.Info("uploaded files", "duration", time.Since(time.UnixMicro(startOSUploadTime)))
	endOsUploadTime := time.Now().UnixMicro()
	s.osDuration.Record(ctx, endOsUploadTime-startOSUploadTime)
	if len(vfs) > 0 {
		if err := s.photoRepo.BulkCreate(ctx, vfs); err != nil {
			cleanup()
			return nil, err
		}
		// Collect created photo IDs for potential cleanup
		for _, photo := range vfs {
			createdPhotoIDs = append(createdPhotoIDs, photo.ID)
		}
	}
	vfs = nil
	createdRecord.VehiclePhotos = bts

	sherr := s.sendToShardari(ctx, createdRecord)
	if sherr == nil || errors.Is(sherr, ErrFailedUpdate) {
		lg.Info("record created", "recordID", record.RecordID)
		return createdRecord, nil
	}
	lg.Warn("failed to send to shardari retrying...", "error", sherr.Error(), "recordID", createdRecord.RecordID)
	if err := s.sendToRetryQueue(ctx, createdRecord, lg); err != nil {
		return nil, err
	}

	lg.Info("record sent for retrying", "recordID", createdRecord.RecordID)
	endProcessDuration := time.Now().UnixMicro()
	s.osDuration.Record(ctx, endProcessDuration-startProcessDuration)
	return createdRecord, nil
}

func (s *VehicleRecordService) Retry(ctx context.Context, createdRecord *entity.VehicleRecord) error {
	lg := s.logger.With("method", "Retry")
	// check if its not in the next day and it's not behind the back of time
	now := time.Now()
	if now.Day() < time.Unix(createdRecord.CreatedAt, 0).Day() {
		return nil
	}
	if now.Unix() < createdRecord.BackoffTime {
		err := s.sendToRetryQueue(ctx, createdRecord, lg)
		return err
	}

	createdRecord.Retries++
	createdRecord.BackoffTime = now.Add(time.Duration(createdRecord.Retries) * (time.Second * 2)).Unix()

	if _, err := s.recordRepo.IncreaseRetry(ctx, createdRecord); err != nil {
		s.logger.Error("failed to increase retry", "error", err.Error())
		return err
	}
	sherr := s.sendToShardari(ctx, createdRecord)
	if sherr == nil || errors.Is(sherr, ErrFailedUpdate) {
		return nil
	}
	if err := s.sendToRetryQueue(ctx, createdRecord, lg); err != nil {
		return err
	}
	lg.Info("record sent for retrying", "recordID", createdRecord.RecordID,
		"backoffTime", createdRecord.BackoffTime, "retry", createdRecord.Retries)
	return nil
}

func (s *VehicleRecordService) sendToRetryQueue(ctx context.Context, createdRecord *entity.VehicleRecord, lg *slog.Logger) error {
	b, err := json.Marshal(&createdRecord)
	if err != nil {
		lg.Error("failed to marshal created record", "error", err)
		return err
	}
	if err = s.rbt.RCH.PublishWithContext(ctx, s.env.RabbitMQInternalExchange, "farin.vehicles.drivers.event.retry",
		false, false, amqp091.Publishing{
			ContentType: "application/json",
			Body:        b,
		}); err != nil {
		lg.Error("failed to send to rabbitmq retry queue", "error", err.Error())
		return err
	}
	return nil
}

func (s *VehicleRecordService) sendToShardari(ctx context.Context, createdRecord *entity.VehicleRecord) error {
	startTime := time.Now().UnixMicro()
	lg := s.logger.With("method", "sendToShardari")
	defer func() {
		lg.Info("send to shardari", "duration", time.Since(time.UnixMicro(startTime)))
	}()

	requestID, plateDetectionID, err := s.tehranSiteRepo.SendVehicleRecord(ctx, createdRecord)
	lg.Info("send to tehran.ir", "requestID", requestID, "plateDetectionID", plateDetectionID)
	if err != nil {
		lg.Error("failed to send vehicle record to tehran.ir", "error", err.Error())
		createdRecord.Sent = false
		if _, err := s.recordRepo.Update(ctx, createdRecord); err != nil {
			lg.Error("failed to update record", "error", err.Error())
			return err
		}
		return err
	} else {
		createdRecord.Sent = true
		createdRecord.TehranRequestID = requestID
		createdRecord.PlateDetectionID = plateDetectionID
		if _, err := s.recordRepo.Update(ctx, createdRecord); err != nil {
			lg.Error("failed to update record", "error", err.Error())
			return ErrFailedUpdate
		}
		endTime := time.Now().UnixMicro()
		s.shardariDuration.Record(ctx, endTime-startTime)
	}

	return nil
}

func (s *VehicleRecordService) FindResend(ctx context.Context, fromDate int64, limit int) error {
	lg := s.logger.With("method", "FindResend")
	records, err := s.recordRepo.GetNotSent(ctx, fromDate, limit)
	if err != nil {
		lg.Error("failed to fetch records", "error", err.Error())
		return nil
	}

	for i, record := range records {
		for k, vf := range record.VehiclePhotos {
			strs := strings.Split(vf.CitizenVehiclePhoto, "/")
			bt, err := downloadDecode(ctx, s, strs, lg)
			if err != nil {
				lg.Warn("failed to download decode for CitizenVehiclePhoto", "error", err.Error())
				break
			}
			records[i].VehiclePhotos[k].CitizenVehiclePhoto = string(bt)

			strs = strings.Split(vf.CitizenVehiclePlateCropPhoto, "/")
			bt, err = downloadDecode(ctx, s, strs, lg)
			if err != nil {
				lg.Warn("failed to download decode for CitizenVehiclePlateCropPhoto", "error", err.Error())
				break
			}
			records[i].VehiclePhotos[k].CitizenVehiclePlateCropPhoto = string(bt)
		}
		if err := s.sendToShardari(ctx, &record); err != nil {
			lg.Warn("failed to send to shardari", "error", err.Error())
			continue
		}
	}
	return nil
}

func (s *VehicleRecordService) Detail(ctx context.Context, id string) (*entity.VehicleRecord, error) {
	lg := s.logger.With("method", "Detail")
	record, err := s.recordRepo.GetByID(ctx, id)
	if err != nil {
		lg.Error("failed to fetch records", "error", err.Error())
		return nil, err
	}

	for k, vf := range record.VehiclePhotos {
		strs := strings.Split(vf.CitizenVehiclePhoto, "/")
		bt, err := downloadDecode(ctx, s, strs, lg)
		if err != nil {
			return nil, err
		}
		record.VehiclePhotos[k].CitizenVehiclePhoto = string(bt)

		strs = strings.Split(vf.CitizenVehiclePlateCropPhoto, "/")
		bt, err = downloadDecode(ctx, s, strs, lg)
		if err != nil {
			return nil, err
		}
		record.VehiclePhotos[k].CitizenVehiclePlateCropPhoto = string(bt)
	}

	return record, nil
}

func downloadDecode(ctx context.Context, s *VehicleRecordService, strs []string, lg *slog.Logger) ([]byte, error) {
	obj, err := s.minio.M.GetObject(ctx, strs[0], strs[1], minio2.GetObjectOptions{})
	if err != nil {
		lg.Error("Failed to get object", "error", err.Error())
		return nil, err
	}
	bf := bytes.NewBuffer([]byte{})
	_, err = io.Copy(bf, obj)
	if err != nil {
		lg.Error("Failed to read object", "error", err.Error())
		return nil, err
	}
	// Calculate the maximum decoded size
	decodedData := bf.Bytes()
	encodedSize := base64.StdEncoding.EncodedLen(len(decodedData))

	// Allocate the destination buffer with the correct size
	encodedData := make([]byte, encodedSize)

	// Decode and get the actual size
	base64.StdEncoding.Encode(encodedData, decodedData)

	return encodedData, nil
}

func plateReplaceAlphabets(input string) (int, error) {
	// Define a map of Persian characters to their corresponding numbers
	replacements := map[string]string{
		"ب":   "02",
		"ج":   "06",
		"ح":   "08",
		"خ":   "09",
		"د":   "10",
		"ذ":   "11",
		"ر":   "12",
		"س":   "15",
		"ص":   "17",
		"ض":   "18",
		"ط":   "19",
		"ظ":   "20",
		"غ":   "22",
		"ق":   "24",
		"ل":   "27",
		"م":   "28",
		"ن":   "29",
		"و":   "30",
		"هـ":  "31",
		"ه":   "31",
		"ی":   "32",
		"الف": "01",
		"ا":   "01",
		"ژ":   "14",
		"چ":   "07",
		"پ":   "03",
		"ت":   "04",
		"ث":   "05",
		"ز":   "13",
		"ش":   "16",
		"ع":   "21",
		"ف":   "23",
		"ک":   "25",
		"گ":   "26",
		"D":   "33",
		"S":   "34",
	}

	result := input

	// For each Persian character in the map, replace it with its corresponding number
	for persianChar, number := range replacements {
		if strings.Contains(result, persianChar) {
			result = strings.ReplaceAll(result, persianChar, "")
			result = result[0:2] + number + result[2:]
		}
	}

	return strconv.Atoi(result)
}
