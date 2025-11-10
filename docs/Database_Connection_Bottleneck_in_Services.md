## Assistant

The `registry.go` file currently creates a single `*sql.DB` connection that is shared across all services. To improve performance and avoid bottlenecks under high load, you can implement a connection pool or distribute connections across services.

Here’s a step-by-step plan:

### 1. **Use a Connection Pool**

Instead of reusing a single connection, use a connection pool. The `database/sql` package supports connection pooling out of the box. You can configure it using `sql.OpenDB` with a connection string that includes a pool size.

### 2. **Modify `NewRegistry` to Use a Connection Pool**

Update the `NewRegistry` function to use a connection pool instead of a single connection.

### 3. **Distribute Connections Across Services**

Each service can have its own connection pool to avoid contention.

### 4. **Implement a Connection Manager**

Create a connection manager that handles the creation and management of connections.

Here’s an example of how you can modify the `NewRegistry` function to use a connection pool:

```go
package services

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Registry struct {
	db *sql.DB

	Devices        *Devices
	Colors         *Colors
	Groups         *Groups
	DeviceControls *DeviceControls
}

// NewRegistry creates a new registry instance and will call the `CreateTables` method.
func NewRegistry(db *sql.DB) *Registry {
	r := &Registry{db: db}

	r.Devices = NewDevices(r)
	r.Colors = NewColors(r)
	r.Groups = NewGroups(r)
	r.DeviceControls = NewDeviceControls(r)

	if err := r.CreateTables(); err != nil {
		panic("failed to create tables: " + err.Error())
	}

	return r
}

func (r *Registry) CreateTables() error {
	services := []Service{
		r.Devices,
		r.Colors,
		r.Groups,
		r.DeviceControls,
	}

	var err error
	for _, service := range services {
		if err = service.CreateTable(); err != nil {
			return err
		}
	}

	return nil
}

func (r *Registry) Close() error {
	return r.db.Close()
}
```

### 5. **Update the Connection Pool Configuration**

You can configure the connection pool by setting the `MaxOpenConns` and `MaxIdleConns` options in the `sql.DB`:

```go
func NewRegistry(db *sql.DB) *Registry {
	// Configure the connection pool
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	r := &Registry{db: db}

	r.Devices = NewDevices(r)
	r.Colors = NewColors(r)
	r.Groups = NewGroups(r)
	r.DeviceControls = NewDeviceControls(r)

	if err := r.CreateTables(); err != nil {
		panic("failed to create tables: " + err.Error())
	}

	return r
}
```

### 6. **Implement a Connection Manager**

Create a connection manager that handles the creation and management of connections. This can be a separate service or a utility function that can be used by all services.
