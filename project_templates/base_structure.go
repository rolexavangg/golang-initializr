package project_templates

import (
	"fmt"
	"strings"
	"text/template"
)

type ProjectConfig struct {
	Name         string
	Dependencies []string
}

func (p *ProjectConfig) GetProjectName() string {
	parts := strings.Split(p.Name, "/")
	return parts[len(parts)-1]
}

func (p *ProjectConfig) HasDependency(name string) bool {
	for _, dep := range p.Dependencies {
		if dep == name {
			return true
		}
	}
	return false
}

func (p *ProjectConfig) BaseFiles() map[string]string {
	files := make(map[string]string)

	readmeTmpl := template.Must(template.New("readme").Parse(readmeTemplate))
	var readmeContent strings.Builder
	readmeTmpl.Execute(&readmeContent, p)
	files["README.md"] = readmeContent.String()

	goModTmpl := template.Must(template.New("gomod").Parse(goModTemplate))
	var goModContent strings.Builder
	goModTmpl.Execute(&goModContent, p)
	files["go.mod"] = goModContent.String()

	dockerfileTmpl := template.Must(template.New("dockerfile").Parse(dockerfileTemplate))
	var dockerfileContent strings.Builder
	dockerfileTmpl.Execute(&dockerfileContent, p)
	files["Dockerfile"] = dockerfileContent.String()

	files[".gitignore"] = gitignoreTemplate

	projectName := p.GetProjectName()
	files["internal/app/app.go"] = fmt.Sprintf(appTemplate, projectName, projectName, projectName, projectName, projectName)

	files["internal/config/config.go"] = configMainTemplate
	files["internal/config/app.go"] = configAppTemplate
	files["internal/config/server.go"] = configServerTemplate
	files["internal/config/utils.go"] = configUtilsTemplate

	files["internal/domain/user.go"] = userDomainTemplate

	files["internal/usecase/user_usecase.go"] = userUsecaseTemplate

	files["internal/repository/user_repository.go"] = userRepositoryTemplate

	files["internal/delivery/http/handler.go"] = httpHandlerTemplate

	files["internal/bootstrap/fx.go"] = bootstrapFxTemplate
	files["internal/bootstrap/logger.go"] = bootstrapLoggerTemplate
	files["internal/bootstrap/http.go"] = bootstrapHttpTemplate

	files[".env.example"] = exampleEnvTemplate
	files["Makefile"] = makefileTemplate

	return files
}

func (p *ProjectConfig) AdditionalFiles() map[string]string {
	files := make(map[string]string)

	if p.HasDependency("http") {
		files["internal/delivery/http/server.go"] = httpServerTemplate
		files["internal/delivery/http/user_handler.go"] = httpUserHandlerTemplate
	}

	if p.HasDependency("grpc") {

		files["internal/config/grpc.go"] = configGRPCTemplate

		files["internal/bootstrap/grpc.go"] = bootstrapGRPCTemplate

		files["internal/delivery/grpc/server.go"] = grpcServerTemplate
		files["internal/delivery/grpc/user_service.go"] = grpcUserServiceTemplate
		files["api/proto/user.proto"] = userProtoTemplate
	}

	if p.HasDependency("postgres") {

		files["internal/config/postgres.go"] = configPostgresTemplate

		files["internal/bootstrap/postgres.go"] = bootstrapPostgresTemplate

		files["internal/repository/postgres/postgres.go"] = postgresTemplate
		files["internal/repository/postgres/user_repository.go"] = postgresUserRepositoryTemplate
	}

	if p.HasDependency("redis") {

		files["internal/config/redis.go"] = configRedisTemplate

		files["internal/bootstrap/redis.go"] = bootstrapRedisTemplate

		files["internal/repository/redis/redis.go"] = redisTemplate
		files["internal/repository/redis/user_cache.go"] = redisUserCacheTemplate
	}

	if p.HasDependency("kafka") {

		files["internal/config/kafka.go"] = configKafkaTemplate

		files["internal/bootstrap/kafka.go"] = bootstrapKafkaTemplate

		files["internal/messaging/kafka/kafka.go"] = kafkaTemplate
		files["internal/messaging/kafka/user_events.go"] = kafkaUserEventsTemplate
	}

	if p.HasDependency("docker") {
		var dockerComposeContent strings.Builder
		dockerComposeTmpl := template.Must(template.New("dockercompose").Parse(dockerComposeTemplate))
		dockerComposeTmpl.Execute(&dockerComposeContent, p)
		files["docker-compose.yml"] = dockerComposeContent.String()
	}

	return files
}

func (p *ProjectConfig) GenerateProject() map[string]string {
	files := p.BaseFiles()

	additionalFiles := p.AdditionalFiles()
	for path, content := range additionalFiles {
		files[path] = content
	}

	return files
}
