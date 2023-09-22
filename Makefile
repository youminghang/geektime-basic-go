# 你可以直接执行 make 命令，也可以单独的命令复制到控制台。
# 注意，如果你是 Windows 并且不是在 WSL 下，
# 要注意文件分隔符使用 Windows 的分隔符。
.PHONY: mock
mock:
	@mockgen -source=./webook/internal/web/jwt/types.go -package=jwtmocks -destination=./webook/internal/web/jwt/mocks/handler.mock.go
	@mockgen -source=./webook/internal/service/user.go -package=svcmocks -destination=./webook/internal/service/mocks/user.mock.go
	@mockgen -source=./webook/internal/service/code.go -package=svcmocks -destination=./webook/internal/service/mocks/code.mock.go
	@mockgen -source=./webook/internal/service/article.go -package=svcmocks -destination=./webook/internal/service/mocks/article.mock.go
	@mockgen -source=./webook/internal/service/sms/types.go -package=smsmocks -destination=./webook/internal/service/sms/mocks/svc.mock.go
	@mockgen -source=./webook/internal/service/oauth2/wechat/types.go -package=wechatmocks -destination=./webook/internal/service/oauth2/wechat/mocks/svc.mock.go
	@mockgen -source=./webook/internal/repository/code.go -package=repomocks -destination=./webook/internal/repository/mocks/code.mock.go
	@mockgen -source=./webook/internal/repository/user.go -package=repomocks -destination=./webook/internal/repository/mocks/user.mock.go
	@mockgen -source=./webook/internal/repository/article_author.go -package=repomocks -destination=./webook/internal/repository/mocks/article_author.mock.go
	@mockgen -source=./webook/internal/repository/article_reader.go -package=repomocks -destination=./webook/internal/repository/mocks/article_reader.mock.go
	@mockgen -source=./webook/internal/repository/dao/user.go -package=daomocks -destination=./webook/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./webook/internal/repository/dao/article/types.go -package=artdaomocks -destination=./webook/internal/repository/dao/article/mocks/article.mock.go
	@mockgen -source=./webook/internal/repository/cache/user.go -package=cachemocks -destination=./webook/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=./webook/pkg/ratelimit/types.go -package=limitmocks -destination=./webook/pkg/ratelimit/mocks/limit.mock.go
	@mockgen -package=redismocks -destination=./webook/internal/repository/cache/redismocks/cmd.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy

.PHONY: e2e
e2e:
	@docker compose -f webook/docker-compose.yaml down
	@docker compose -f webook/docker-compose.yaml up -d
	@go test -race ./webook/... -tags=e2e
	@docker compose -f webook/docker-compose.yaml down
.PHONY: e2e_up
e2e_up:
	@docker compose -f webook/docker-compose.yaml up -d
.PHONY: e2e_down
e2e_down:
	@docker compose -f webook/docker-compose.yaml down