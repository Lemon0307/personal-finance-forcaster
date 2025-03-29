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
        response = []
        if not jsonData:
            return jsonify({"message": "cannot find JSON data provided"}, 400)
        items = jsonData.get("Items")
        
        for item in items:
            # extract transactions from json data
            total_spent = list(map(lambda x: x.get('Amount'), item.get("TotalSpent")))
            total_earned = list(map(lambda x: x.get('Amount'), item.get("TotalEarned")))
            dates = list(map(lambda x: x.get('Date'), item.get("TotalSpent")))

            # sort dates by the most recent date

            print(dates)
            dates.sort(reverse=True)
            recent_month_year = dates[0].split('-')
            recent_month = int(recent_month_year[1])
            recent_year = int(recent_month_year[0])
            # forecast the user's spending
            forecasted_spending, error = forecast(total_spent, int(months), 1, 1, 1)
            if error != None:
                return jsonify({"error": error})
            # forecast the user's earnings
            forecasted_earning, error = forecast(total_earned, int(months), 1, 1, 1)
            print(recent_month)

            res_spent = []
            res_earned = []

            # generate new month and year for forecasted transactions
            for i in range(len(forecasted_spending)):
                month = ((recent_month - 1 + i + 1) % 12) + 1
                year = recent_year + ((recent_month - 1 + i + 1) // 12)
                print(str(year) + '-' + str(month))
                res_spent.append({"Date": str(year) + '-' + str(month).zfill(2), "Amount": forecasted_spending[i]})
                res_earned.append({"Date": str(year) + '-' + str(month).zfill(2), "Amount": forecasted_earning[i]})
            
            # prepare data and append to the response array
            sub_response = {}
            sub_response['item_name'] = item.get("ItemName")
            sub_response["total_spending"] = item.get("TotalSpent")
            sub_response["total_earning"] = item.get("TotalEarned")
            sub_response["forecasted_spending"] = res_spent
            sub_response["forecasted_earning"] = res_earned

            sub_response["net_cash_flow"] = round(float(sum(forecasted_earning) - sum(forecasted_spending)), 2)
            sub_response["recommended_budget"] = round(float(mean_value(total_spent + forecasted_spending)), 2)
            response.append(sub_response)

        return jsonify(response)

api.add_resource(Forecast, '/forecast')

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)