docker run -it --rm --gpus all -v /home/juwon/triton_vitpose/:/workspace 105e6f8b8954 trtexec --onnx=/workspace/vitpose-b-multi-coco.onnx --fp16 --minShapes='input':1x3x256x192 --optShapes='input':8x3x256x192 --maxShapes='input':16x3x256x192 --saveEngine=/workspace/model.plan
mkdir -p pose_model_zoo/vitpose/1
mv model.plan pose_model_zoo/vitpose/1
docker build -t triton_vitpose .
docker run --rm --gpus all -p 8000:8000 -p 8001:8001 -p 8002:8002 -v $pwd/pose_model_zoo:/models triton_vitpose tritonserver --model-repository=/models