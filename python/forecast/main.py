from fastapi import FastAPI
from ARIMA import forecast_data
app = FastAPI()

@app.get("/forecast/get_data")
async def get_data():
    pass

@app.get("/forecast/predict")
async def predict():
    pass

@app.get("/forecast/recommend")
async def recommend():
    pass