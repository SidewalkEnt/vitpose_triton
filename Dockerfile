FROM nvcr.io/nvidia/tritonserver:24.05-py3
RUN pip install opencv-python
RUN apt-get update -y
RUN apt-get install -y libgl1-mesa-glx
COPY ./pose_model_zoo /models

CMD ["tritonserver", "--model-repository=/models"]