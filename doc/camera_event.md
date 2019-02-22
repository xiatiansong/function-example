# 摄像头事件触发规则引擎

## 文件上传完成事件
规则引擎输入参数：
```
{
    "token": "xxxx", // 请求token，调用miot API时需要
    "data":{
        "did": "12941756",
        "uid": 49057842,
        "model": "mijia.camera.v1",
        "pd_id": 999,
        "attr": "event.third_part_file_upload",
        "status": {
            "thirdpartSegment":{
                "duration":10,
                "expireTime":1533265848969,
                "offset":0,
                "createTime":1533265748969,
                "videoStoreId":"videoStoreId",
                "userId":1000000,
                "imgStoreId":"imgStoreId",
                "did":"did",
                "fileType":"VIDEO",
                "videoSize":10,
                "fileId":"7e9dddea-07dc-48c2-ae91-6c015b0f4a11",
                "extraInfo":"thirdPartSegment_extraInfo"
            },
            "extraInfo":"extraInfo"
        }
    }
}

其中thirdpartSegment 包含以下成员变量:

struct ThirdPartSegment {
    1: required i64 userId;
    2: required string did;
    3: required string fileId;
    4: optional i32 offset;
    5: optional string location;
    6: optional string fileType;
    7: optional i64 createTime;
    8: optional i64 expireTime;
    9: optional i64 videoSize;
    10: optional double duration;
    11: optional string videoStoreId;
    12: optional string imgStoreId;
    13: optional string extraInfo;
}
```

## 人脸识别完成事件
规则引擎输入参数：
```
{
    "token": "xxxx", // 请求token，调用miot API时需要
    "data":{
        "did": "12941756",
        "uid": 49057842,
        "model": "mijia.camera.v1",
        "pd_id": 999,
        "attr": "event.camera_face_detection",
        "status": {
            "faceInfoMetas":[
                {
                "score":2000,
                "ageConfidence":0.9555555,
                "genderConfidence":0.9,
                "gender":1,
                "wearGlasses":true,
                "faceIdStr":"11782115086565888",
                "age":23
                },
            ]
        }
    }
}
其中faceInfoMetas包含的元素，可以解析为:
struct FaceInfoMeta {
    1: required string faceIdStr;
    2: optional double score;
    3: optional i32 age;
    4: optional double ageConfidence;
    5: optional i32 gender;
    6: optional double genderConfidence;
    7: optional bool wearGlasses;
}
```