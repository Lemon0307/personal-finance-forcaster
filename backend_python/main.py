from flask import Flask, request
from flask_restful import Api, Resource

app = Flask(__name__)
api = Api(app)

class Forecast(Resource):
    def get(self): 
        pass
    def post(self):
        data = request.get_json()
        

api.add_resource(Forecast, '/forecast')

if __name__ == '__main__':
    app.run(debug=True)