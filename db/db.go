package db

import (
	"encoding/binary"
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

var (
	db         *bolt.DB
	taskBucket = []byte("tasks")
)

type Task struct {
	Key       int       `json:"key"`
	Value     string    `json:"value"`
	TimeAdded time.Time `json:"timeAdded"`
	Completed bool      `json:"completed"`
}

func (t Task) FilterValue() string {
	return t.Value
}

func Open(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return err
	}

	return db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists(taskBucket)
		if err != nil {
			return err
		}
		return nil
	})
}

func CreateTask(t string) (Task, error) {
	var task Task
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id := int(id64)
		task = Task{
			Key:       id,
			Value:     t,
			TimeAdded: time.Now(),
			Completed: false,
		}

		value, err := json.Marshal(task)
		if err != nil {
			return err
		}

		return b.Put(itob(id), value)

	})
	if err != nil {
		return task, err
	}
	return task, nil
}

type status int

const (
	active status = iota
	done
	all
)

func list(s status) ([]Task, error) {
	var tasks []Task
	err := db.View(func(t *bolt.Tx) (err error) {
		b := t.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			task := Task{}
			err = json.Unmarshal(v, &task)
			if err != nil {
				return err
			}

			switch s {
			case active:
				if task.Completed != true {
					tasks = append(tasks, task)
				}
			case done:
				if task.Completed == true {
					tasks = append(tasks, task)
				}
			case all:
				tasks = append(tasks, task)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func AllTasks() ([]Task, error) {
	return list(all)
}

func CompletedTasks() ([]Task, error) {
	return list(done)
}

func ActiveTasks() ([]Task, error) {
	return list(active)
}

func updateStatus(key int, s status) (Task, error) {
	task := Task{}
	return task, db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(taskBucket)
		id := itob(key)
		value := b.Get(id)

		if value != nil {
			err := json.Unmarshal(value, &task)
			if err != nil {
				return err
			}

			switch s {
			case active:
				task.Completed = false
			case done:
				task.Completed = true
			}

			buf, err := json.Marshal(task)
			if err != nil {
				return err
			}

			err = b.Put(id, buf)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func MarkDone(key int) (Task, error) {
	return updateStatus(key, done)
}

func MarkActive(key int) (Task, error) {
	return updateStatus(key, active)
}

func itob(i int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
