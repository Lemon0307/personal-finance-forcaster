from fastapi import FastAPI, Path

app = FastAPI()

@app.get("/")