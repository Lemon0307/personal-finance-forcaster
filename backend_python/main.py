from flask import Flask, request, jsonify
from flask_restful import Api, Resource
from forecast import forecast, mean_value
import numpy as np
from flask_cors import CORS, cross_origin

app = Flask(__name__)

api = Api(app)
CORS(app, supports_credentials=True)

class Forecast(Resource):
    @cross_origin()
    def post(self):
        jsonData = request.get_json()
        # get the number of months to forecast from url
        months = request.args.get('months')
        if not jsonData:
            return jsonify({"message": "cannot find JSON data provided"}, 400)
        items = jsonData.get("Items")
        for item in items:
            transactions = item.get("Transactions")
            forecasted_transactions = forecast(transactions, int(months), 12, 1, 12)
            print(forecasted_transactions)

api.add_resource(Forecast, '/forecast')

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)