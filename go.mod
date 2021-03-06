module github.com/rai-project/dlframework

go 1.12

require (
	cloud.google.com/go v0.34.0
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78
	github.com/GeertJohan/go-sourcepath v0.0.0-20150925135350-83e8b8723a9b
	github.com/Masterminds/semver v1.4.2
	github.com/Microsoft/go-winio v0.4.11
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5
	github.com/PuerkitoBio/purell v1.1.0
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578
	github.com/Shopify/sarama v1.19.0
	github.com/StackExchange/wmi v0.0.0-20180116203802-5d049714c4a6
	github.com/Unknwon/com v0.0.0-20151008135407-28b053d5a292
	github.com/VividCortex/ewma v0.0.0-20170804035156-43880d236f69
	github.com/VividCortex/robustly v0.0.0-20180323182711-6d4b5835c602
	github.com/airbrake/gobrake v3.6.1+incompatible
	github.com/anthonynsimon/bild v0.10.0
	github.com/apache/thrift v0.0.0-20181207211846-208a048dc440
	github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6
	github.com/asaskevich/govalidator v0.0.0-20180315120708-ccb8e960c48f
	github.com/aws/aws-sdk-go v1.16.2
	github.com/bamiaux/rez v0.0.0-20170731184118-29f4463c688b
	github.com/benesch/cgosymbolizer v0.0.0-20180702220239-70e1ee2b39d3
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d
	github.com/boltdb/bolt v1.3.1
	github.com/bugsnag/osext v0.0.0-20130617224835-0dd3f918b21b
	github.com/carlescere/scheduler v0.0.0-20170109141437-ee74d2f83d82
	github.com/cenkalti/backoff v2.1.0+incompatible
	github.com/cockroachdb/cmux v0.0.0-20170110192607-30d10be49292
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd
	github.com/containerd/continuity v0.0.0-20181203112020-004b46473808
	github.com/coreos/etcd v3.3.10+incompatible
	github.com/coreos/go-semver v0.2.0
	github.com/coreos/go-systemd v0.0.0-20181031085051-9002847aa142
	github.com/cpuguy83/go-md2man v1.0.8
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/distribution v2.7.0+incompatible
	github.com/docker/docker v0.0.0-20181207101903-a4a816b6bbed
	github.com/docker/go-connections v0.0.0-20180821093606-97c2040d34df
	github.com/docker/go-plugins-helpers v0.0.0-20181025120712-1e6269c305b8
	github.com/docker/go-units v0.3.3
	github.com/dustin/go-humanize v1.0.0
	github.com/eapache/go-resiliency v1.1.0
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21
	github.com/eapache/queue v1.1.0
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/evalphobia/logrus_fluent v0.4.0
	github.com/evalphobia/logrus_kinesis v0.2.0
	github.com/facebookgo/freeport v0.0.0-20150612182905-d4adf43b75b9
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052
	github.com/fatih/color v1.7.0
	github.com/fatih/structs v1.1.0
	github.com/fluent/fluent-logger-golang v1.3.0
	github.com/flyaways/golang-lru v0.0.0-20180528053744-bdd6594a7c32
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gbbr/memstats v0.0.0-20141117200020-3f4151ce3189
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-logfmt/logfmt v0.4.0
	github.com/go-ole/go-ole v1.2.1
	github.com/go-openapi/analysis v0.17.2
	github.com/go-openapi/errors v0.19.0
	github.com/go-openapi/jsonpointer v0.17.2
	github.com/go-openapi/jsonreference v0.17.2
	github.com/go-openapi/loads v0.19.0
	github.com/go-openapi/runtime v0.18.0
	github.com/go-openapi/spec v0.18.0
	github.com/go-openapi/strfmt v0.18.0
	github.com/go-openapi/swag v0.18.0
	github.com/go-openapi/validate v0.18.0
	github.com/go-playground/locales v0.12.1
	github.com/go-playground/universal-translator v0.16.0
	github.com/gogo/googleapis v1.1.0
	github.com/gogo/protobuf v1.1.1
	github.com/golang/protobuf v0.0.0-20181128192352-1d3f30b51784
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db
	github.com/gonum/blas v0.0.0-20180125090452-e7c5890b24cf
	github.com/gonum/internal v0.0.0-20181124074243-f884aa714029
	github.com/google/go-querystring v1.0.0
	github.com/google/gops v0.3.5
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/grpc-gateway v1.4.1
	github.com/hashicorp/consul v1.4.0
	github.com/hashicorp/go-cleanhttp v0.5.0
	github.com/hashicorp/go-getter v0.0.0-20181119194526-bd1edc22f8ea
	github.com/hashicorp/go-rootcerts v0.0.0-20160503143440-6bb64b370b90
	github.com/hashicorp/go-safetemp v1.0.0
	github.com/hashicorp/go-version v1.0.0
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/serf v0.8.1
	github.com/iancoleman/orderedmap v0.0.0-20181121102841-22c6ecc9fe13
	github.com/ianlancetaylor/cgosymbolizer v0.0.0-20170921033129-f5072df9c550
	github.com/ianlancetaylor/demangle v0.0.0-20181102032728-5e5cf60278f6
	github.com/imdario/mergo v0.3.6
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/intel-go/cpuid v0.0.0-20181003105527-1a4a6f06a1c6
	github.com/jessevdk/go-flags v1.4.0
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af
	github.com/json-iterator/go v1.1.5
	github.com/junegunn/go-shellwords v0.0.0-20170411071455-02e3cf038dce
	github.com/k0kubun/pp v2.3.0+incompatible
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1
	github.com/klauspost/shutdown2 v1.1.0
	github.com/knq/jwt v0.0.0-20180925223530-fc44a4704737
	github.com/knq/pemutil v0.0.0-20180607233853-a6a7785bc45a
	github.com/knq/sdhook v0.0.0-20181029224735-f9c6ae1bf3ae
	github.com/konsorten/go-windows-terminal-sequences v1.0.1
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515
	github.com/labstack/echo v0.0.0-20181123063703-c7eb8da9ec73
	github.com/labstack/gommon v0.2.8
	github.com/levigross/grequests v0.0.0-20160921031216-3f92c0acb6cd
	github.com/magiconair/properties v1.8.0
	github.com/mailru/easyjson v0.0.0-20180823135443-60711f1a8329
	github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-isatty v0.0.4
	github.com/mattn/go-runewidth v0.0.3
	github.com/mitchellh/colorstring v0.0.0-20150917214807-8631ce90f286
	github.com/mitchellh/go-homedir v1.0.0
	github.com/mitchellh/go-testing-interface v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742
	github.com/nicolai86/instruments v0.0.0-20170630130909-a667d8f6e278
	github.com/olekukonko/tablewriter v0.0.1
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/opencontainers/runc v0.0.0-20181204161456-25f3f893c86d
	github.com/opentracing-contrib/go-observer v0.0.0-20170622124052-a52f23424492
	github.com/opentracing-contrib/perfevents v0.0.0-20171011010702-a7a7e747782c
	github.com/opentracing/opentracing-go v1.0.2
	github.com/openzipkin/zipkin-go-opentracing v0.3.4
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pelletier/go-toml v1.2.0
	github.com/philhofer/fwd v1.0.0
	github.com/pierrec/lz4 v0.0.0-20181005164709-635575b42742
	github.com/pkg/errors v0.8.0
	github.com/pmezard/go-difflib v1.0.0
	github.com/rai-project/acl v0.0.0-20181119122707-037e0eb4d746
	github.com/rai-project/aws v0.0.0-20181119122706-0989b18a4aeb
	github.com/rai-project/batching v0.0.0-20181120145124-f672e97e4157
	github.com/rai-project/cmd v0.0.0-20181119122707-69f8e596de1c
	github.com/rai-project/config v0.0.0-20181119122707-4d0de45fe6c1
	github.com/rai-project/cpu v0.0.0-20181119122707-5bc145cba8e0
	github.com/rai-project/database v0.0.0-20181202070225-ef4eddec538e
	github.com/rai-project/dldataset v0.0.0-20181119123731-c28e4f89153b
	github.com/rai-project/docker v0.0.0-20181119123731-86779e24c596
	github.com/rai-project/downloadmanager v0.0.0-20181119123731-abede5faa82f
	github.com/rai-project/evaluation v0.0.0-20181204141616-e0b15a6a3390
	github.com/rai-project/go-cupti v0.0.0-20181121031418-ffa685874f19
	github.com/rai-project/go-libjpeg v0.0.0-20181119123732-40fd7b1bcbeb
	github.com/rai-project/godotenv v0.0.0-20180908223441-72ca456a35f4
	github.com/rai-project/googlecloud v0.0.0-20181119123731-8938dc83da61
	github.com/rai-project/grpc v0.0.0-20181121055653-de384a740c84
	github.com/rai-project/image v0.0.0-20181130151553-a6ebd1b24d2f
	github.com/rai-project/ldcache v0.0.0-20181119123732-af85cb316a45
	github.com/rai-project/libkv v0.0.0-20181119123731-49a14e78d856
	github.com/rai-project/lock v0.0.0-20181119123731-67e734de309b
	github.com/rai-project/logger v0.0.0-20181119115247-3edfaed4af1c
	github.com/rai-project/machine v0.0.0-20181119123731-e20c9a49017c
	github.com/rai-project/model v0.0.0-20181119123731-66be2e1deaae
	github.com/rai-project/monitoring v0.0.0-20171102205001-498a6c4f3e85
	github.com/rai-project/nvidia-smi v0.0.0-20181121005638-5f6bbd426877
	github.com/rai-project/nvml-go v0.0.0-20181121025807-4e1189bcc320
	github.com/rai-project/parallel v0.0.0-20181119123731-16f7855030a4
	github.com/rai-project/passlib v0.0.0-20181013114510-cdf39dc0b8ea
	github.com/rai-project/pipeline v0.0.0-20181119123731-78ece549980a
	github.com/rai-project/registry v0.0.0-20181119122707-e1a4540d29a8
	github.com/rai-project/serializer v0.0.0-20181119122706-c73b52fef201
	github.com/rai-project/synthetic_load v0.0.0-20181202050407-58faa17e0d50
	github.com/rai-project/tegra v0.0.0-20181119122707-1d9901ca382b
	github.com/rai-project/tracer v0.0.0-20181129032001-41b5c08e8b8e
	github.com/rai-project/utils v0.0.0-20181119122706-be23e9dad62b
	github.com/rai-project/uuid v0.0.0-20181119122706-2a4c8b922cc6
	github.com/rai-project/vipertags v0.0.0-20181119122706-8cbaab517f5d
	github.com/rai-project/web v0.0.0-20181111061007-97ba6d7fd665
	github.com/rainycape/dl v0.0.0-20151222075243-1b01514224a1
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a
	github.com/russross/blackfriday v1.5.2
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec
	github.com/schollz/progressbar v0.0.0-20181102035236-74139f27599a
	github.com/sebest/logrusly v0.0.0-20180315190218-3235eccb8edc
	github.com/seehuhn/mt19937 v0.0.0-20180715112136-cc7708819361
	github.com/segmentio/go-loggly v0.5.0
	github.com/shirou/gopsutil v2.18.11+incompatible
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/afero v1.1.2
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v0.0.0-20181127133106-d2d81d9a96e2
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.1
	github.com/stretchr/objx v0.1.1
	github.com/stretchr/testify v1.2.2
	github.com/tinylib/msgp v1.0.2
	github.com/uber/jaeger v1.8.2
	github.com/uber/jaeger-client-go v2.15.0+incompatible
	github.com/uber/jaeger-lib v1.5.0
	github.com/ugorji/go v1.1.1
	github.com/ulikunitz/xz v0.5.5
	github.com/ulule/deepcopier v0.0.0-20171107155558-ca99b135e50f
	github.com/unixpickle/anydiff v0.0.0-20170906235721-55bea27bf48b
	github.com/unixpickle/anynet v0.0.0-20170909172929-016782221a5a
	github.com/unixpickle/anyvec v0.0.0-20170908190750-59aa66ba0472
	github.com/unixpickle/autofunc v0.0.0-20170112172612-f27a3f82164a
	github.com/unixpickle/essentials v0.0.0-20180916162721-ae02bc395f1d
	github.com/unixpickle/mnist v0.0.0-20170128023510-751c4271cf3a
	github.com/unixpickle/num-analysis v0.0.0-20161229165253-c45203c63047
	github.com/unixpickle/serializer v0.0.0-20170723202158-c6c092dc55bb
	github.com/unixpickle/sgd v0.0.0-20161225162810-0e3d4c9d317b
	github.com/unixpickle/tensor v0.0.0-20170114180418-7295881ed12b
	github.com/unixpickle/weakai v0.0.0-20170623211141-247102c87396
	github.com/valyala/bytebufferpool v1.0.0
	github.com/valyala/fasttemplate v0.0.0-20170224212429-dcecefd839c4
	github.com/wercker/journalhook v0.0.0-20180428041537-5d0a5ae867b3
	github.com/xordataexchange/crypt v0.0.0-20170626215501-b2862e3d0a77
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9
	golang.org/x/image v0.0.0-20181116024801-cd38e8056d9b
	golang.org/x/net v0.0.0-20181207154023-610586996380
	golang.org/x/oauth2 v0.0.0-20181203162652-d668ce993890
	golang.org/x/sync v0.0.0-20181108010431-42b317875d0f
	golang.org/x/sys v0.0.0-20181206074257-70b957f3b65e
	golang.org/x/text v0.3.0
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c
	google.golang.org/api v0.0.0-20181206211257-1a5ef82f9af4
	google.golang.org/appengine v1.3.0
	google.golang.org/genproto v0.0.0-20181202183823-bd91e49a0898
	google.golang.org/grpc v1.17.0
	gopkg.in/VividCortex/ewma.v1 v1.1.1
	gopkg.in/cheggaaa/pb.v2 v2.0.6
	gopkg.in/fatih/color.v1 v1.7.0
	gopkg.in/gemnasium/logrus-airbrake-hook.v3 v3.0.2
	gopkg.in/gemnasium/logrus-graylog-hook.v2 v2.0.7
	gopkg.in/go-playground/validator.v9 v9.23.0
	gopkg.in/hlandau/easymetric.v1 v1.0.0
	gopkg.in/hlandau/measurable.v1 v1.0.1
	gopkg.in/mattn/go-colorable.v0 v0.0.9
	gopkg.in/mattn/go-isatty.v0 v0.0.4
	gopkg.in/mattn/go-runewidth.v0 v0.0.3
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	gopkg.in/yaml.v2 v2.2.2
	upper.io/db.v3 v3.5.5+incompatible
)
