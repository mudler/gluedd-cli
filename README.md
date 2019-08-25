# gluedd-cli

Usage example of [gluedd](https://github.com/mudler/gluedd/), inspired from [livedetect](https://github.com/jolibrain/livedetect).

It reads JPEG streams and webcam and passes to Deepdetect. It creates a webserver that can be used to debug the models. It also allows to update Items of an OpenHab instance, allowing to connect DeepDetect models with 

## Requirements

To follow this example case you need (on a rpi3):

- Docker (to run the DeepDetect API)
- Golang (to build this project)

## Run:

### 1) Setup [Deepdetect](https://www.deepdetect.com/) API server:

(See also the [livedetect](https://github.com/jolibrain/livedetect/wiki/Step-by-step-for-Raspberry-Pi-3) step-by-step guide)

Starts the [Deepdetect](https://www.deepdetect.com/) API (on a rpi3 with ncnn in this example): 

    $> docker run --restart=always --name deepdetect -d -p 8890:8080 -v $HOME/models:/opt/models jolibrain/deepdetect_ncnn_pi3

Fixup model folder permissions:

    $> sudo chown -R $(id -u ${USER}):$(id -g ${USER}) $HOME/models 

Create a service:

    $> curl -X PUT http://localhost:8890/services/squeezenet_ssd_voc -d '{
        "description": "Squeezenet SSD",
        "model": {
            "repository": "/opt/models/squeezenet_ssd_voc",
            "create_repository": true,
            "init":"https://www.deepdetect.com/models/init/ncnn/squeezenet_ssd_voc_ncnn_300x300.tar.gz"
        },
        "mllib": "ncnn",
        "type": "supervised",
        "parameters": {
            "input": {
                "connector": "image"
            }
        }
    }'

### 2) Build the project (you need go)

    $> go get github.com/mudler/gluedd-cli
    $> pushd $GOPATH/src/github.com/mudler/gluedd-cli
    $> make build

### 3) Create config file

    # Deep detect configs
    api_server: "http://localhost:8080/" # Deep detect api
    service: "squeezenet_ssd_voc"
    buffer_size: 1 # 0 to disable fixed-buffering
    preview: true
    confidence: 0.3

    # Stream
    base_url: "http://0.0.0.0:4000/" # live preview stream url
    stream_url: "http://192.168.1.2:88/cgi-bin/CGIProxy.fcgi?cmd=snapPicture2&usr=&pwd=" # Remote Jpeg stream

    # Openhab
    vehicle_item: "Camera1VehicleDetection" 
    person_item: "Camera1HumanDetection"
    animal_item: "Camera1AnimalDetection"
    openhab_url: "http://openhabianpi:8080"

    # Resize options

    ## Webcam
    webcam_width: 300
    webcam_height: 300

    ## Jpeg stream resize
    image_width: 300
    image_height: 300

    ## General resizing options
    resize: true # Enable image resizing
    approx: true # Resize approximation

### 4) Run it

Against a local webcam:

    $> ./gluedd-cli --config config.yaml webcam 0

Predict from a JPEG stream:

    $> ./gluedd-cli --config config.yaml stream

Predict from a JPEG stream and send Item updates to openhab:

    $> ./gluedd-cli --config config.yaml openhab

If you configured ```preview: true``` you can point your browser to ```base_url``` to see the prediction live stream. (e.g. http://localhost:4000/ in the following example)
