docker build -t triton_vitpose .
docker run --rm --gpus all -p 8000:8000 -p 8001:8001 -p 8002:8002 -v $pwd/pose_model_zoo:/models triton_vitpose tritonserver --model-repository=/models