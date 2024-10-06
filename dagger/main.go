// A generated module for DaggerHelloWorld functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/dagger-hello-world/internal/dagger"
	"fmt"
	"math"
	"math/rand"
)

type DaggerHelloWorld struct{}

// Returns a container that echoes whatever string argument is provided
func (m *DaggerHelloWorld) BuildEnv(source *dagger.Directory) *dagger.Container {
	mavenCache := dag.CacheVolume("maven")

	return dag.Container().From("maven:3.9.9").WithDirectory("/src", source.WithoutDirectory("dagger")).WithMountedCache("/root/.m2", mavenCache).WithWorkdir("/src")
}

func (m *DaggerHelloWorld) Build(ctx context.Context, source *dagger.Directory) *dagger.File {
	return m.BuildEnv(source).WithExec([]string{"mvn", "-B", "-DskipTests", "clean", "package"}).File("target/my-app-1.0-SNAPSHOT.jar")
}

func (m *DaggerHelloWorld) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.BuildEnv(source).WithExec([]string{"mvn", "test"}).Stdout(ctx)
}

func (m *DaggerHelloWorld) Publish(ctx context.Context, source *dagger.Directory) (string, error) {
	// call Dagger Function to run unit tests
	_, err := m.Test(ctx, source)
	if err != nil {
		return "", err
	}
	// call Dagger Function to build the application image
	// publish the image to ttl.sh
	return dag.Container().
		From("eclipse-temurin:17-alpine").
		WithFile("/app/my-app-1.0-snapshot", m.Build(ctx, source)).
		WithEntrypoint([]string{"java", "-jar", "/app/my-app-1.0-snapshot"}).
		Publish(ctx, fmt.Sprintf("ttl.sh/dagger-maven-%.0f", math.Floor(rand.Float64()*10000000))) //#nosec
}
