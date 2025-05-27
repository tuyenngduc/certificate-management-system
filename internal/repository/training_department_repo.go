package repository

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrainingDepartmentRepository struct {
	FacultyCol  *mongo.Collection
	ClassCol    *mongo.Collection
	LecturerCol *mongo.Collection
}

func (r *TrainingDepartmentRepository) GetFacultyByCode(ctx context.Context, code string) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.FacultyCol.FindOne(ctx, bson.M{"code": code}).Decode(&faculty)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &faculty, err
}
func NewTrainingDepartmentRepository(db *mongo.Database) *TrainingDepartmentRepository {
	return &TrainingDepartmentRepository{
		FacultyCol:  db.Collection("faculties"),
		ClassCol:    db.Collection("classes"),
		LecturerCol: db.Collection("lecturers"),
	}
}
func (r *TrainingDepartmentRepository) EnsureFacultyCodeUniqueIndex(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"code": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := r.FacultyCol.Indexes().CreateOne(ctx, indexModel)
	return err
}

// Faculty CRUD
func (r *TrainingDepartmentRepository) CreateFaculty(ctx context.Context, faculty *models.Faculty) error {
	_, err := r.FacultyCol.InsertOne(ctx, faculty)
	return err
}
func (r *TrainingDepartmentRepository) GetAllFaculties(ctx context.Context) ([]models.Faculty, error) {
	var faculties []models.Faculty
	cursor, err := r.FacultyCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var f models.Faculty
		if err := cursor.Decode(&f); err == nil {
			faculties = append(faculties, f)
		}
	}
	return faculties, nil
}
func (r *TrainingDepartmentRepository) GetFacultyByID(ctx context.Context, id primitive.ObjectID) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.FacultyCol.FindOne(ctx, bson.M{"_id": id}).Decode(&faculty)
	if err != nil {
		return nil, err
	}
	return &faculty, nil
}

func (r *TrainingDepartmentRepository) UpdateFaculty(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.FacultyCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}
func (r *TrainingDepartmentRepository) DeleteFaculty(ctx context.Context, id primitive.ObjectID) (bool, error) {
	res, err := r.FacultyCol.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
func (r *TrainingDepartmentRepository) FindFacultyByCode(ctx context.Context, code string) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.FacultyCol.FindOne(ctx, bson.M{"code": code}).Decode(&faculty)
	if err != nil {
		return nil, err
	}
	return &faculty, nil
}

// Class CRUD
func (r *TrainingDepartmentRepository) GetClassesByFacultyID(ctx context.Context, facultyID primitive.ObjectID) ([]models.Class, error) {
	var classes []models.Class
	cursor, err := r.ClassCol.Find(ctx, bson.M{"faculty_id": facultyID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var class models.Class
		if err := cursor.Decode(&class); err == nil {
			classes = append(classes, class)
		}
	}
	return classes, nil
}
func (r *TrainingDepartmentRepository) FindClassByCode(ctx context.Context, code string) (*models.Class, error) {
	var class models.Class
	err := r.ClassCol.FindOne(ctx, bson.M{"code": code}).Decode(&class)
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *TrainingDepartmentRepository) CreateClass(ctx context.Context, class *models.Class) error {
	_, err := r.ClassCol.InsertOne(ctx, class)
	return err
}
func (r *TrainingDepartmentRepository) GetAllClasses(ctx context.Context) ([]models.Class, error) {
	var classes []models.Class
	cursor, err := r.ClassCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var c models.Class
		if err := cursor.Decode(&c); err == nil {
			classes = append(classes, c)
		}
	}
	return classes, nil
}

func (r *TrainingDepartmentRepository) UpdateClass(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.ClassCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *TrainingDepartmentRepository) GetClassByID(ctx context.Context, id primitive.ObjectID) (*models.Class, error) {
	var class models.Class
	err := r.ClassCol.FindOne(ctx, bson.M{"_id": id}).Decode(&class)
	if err != nil {
		return nil, err
	}
	return &class, nil
}
func (r *TrainingDepartmentRepository) DeleteClass(ctx context.Context, id primitive.ObjectID) (bool, error) {
	res, err := r.ClassCol.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}

// Lecturer CRUD
func (r *TrainingDepartmentRepository) CreateLecturer(ctx context.Context, lecturer *models.Lecturer) error {
	_, err := r.LecturerCol.InsertOne(ctx, lecturer)
	return err
}
func (r *TrainingDepartmentRepository) GetAllLecturers(ctx context.Context) ([]models.Lecturer, error) {
	var lecturers []models.Lecturer
	cursor, err := r.LecturerCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var l models.Lecturer
		if err := cursor.Decode(&l); err == nil {
			lecturers = append(lecturers, l)
		}
	}
	return lecturers, nil
}
func (r *TrainingDepartmentRepository) GetLecturersByFacultyID(ctx context.Context, facultyID primitive.ObjectID) ([]models.Lecturer, error) {
	var lecturers []models.Lecturer
	cursor, err := r.LecturerCol.Find(ctx, bson.M{"faculty_id": facultyID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var lecturer models.Lecturer
		if err := cursor.Decode(&lecturer); err == nil {
			lecturers = append(lecturers, lecturer)
		}
	}
	return lecturers, nil
}
func (r *TrainingDepartmentRepository) GetLecturerByID(ctx context.Context, id primitive.ObjectID) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	err := r.LecturerCol.FindOne(ctx, bson.M{"_id": id}).Decode(&lecturer)
	if err != nil {
		return nil, err
	}
	return &lecturer, nil
}
func (r *TrainingDepartmentRepository) UpdateLecturer(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.LecturerCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}
func (r *TrainingDepartmentRepository) DeleteLecturer(ctx context.Context, id primitive.ObjectID) (bool, error) {
	res, err := r.LecturerCol.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}

func (r *TrainingDepartmentRepository) FindLecturerByCode(ctx context.Context, code string) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	err := r.LecturerCol.FindOne(ctx, bson.M{"code": code}).Decode(&lecturer)
	if err != nil {
		return nil, err
	}
	return &lecturer, nil
}
