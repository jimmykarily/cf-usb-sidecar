FROM cf-usb-sidecar-buildbase:latest as builder

COPY . /go/src/github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-mssql
ENV GOPATH /go
ENV GOBIN /go/bin
WORKDIR /go/src/github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-mssql

RUN make tools
RUN make build
RUN chmod +x SIDECAR_HOME/bin/run.sh

FROM cf-usb-sidecar:latest
COPY --from=builder /go/src/github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-mssql/SIDECAR_HOME/ /catalog-service-manager/

ENTRYPOINT ["sh", "/catalog-service-manager/bin/run.sh"]
