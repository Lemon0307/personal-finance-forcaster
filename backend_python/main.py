from flask import Flask, request, jsonify
from flask_restful import Api, Resource
from forecast import forecast, mean_value, knapsack
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
        total_budget = sum(map(lambda x: x.get('BudgetCost'), items))
        for item in items:
            total_spent = list(map(lambda x: x.get('Amount'), item.get("TotalSpent")))
            total_earned = list(map(lambda x: x.get('Amount'), item.get("TotalEarned")))
            dates = list(map(lambda x: {'month': x.get('Month'), 
                                        'year': x.get('Year')}, item.get("TotalSpent")))
            print(total_spent)
            forecasted_spending, error = forecast(total_spent, int(months), 1, 1, 1)
            print(forecasted_spending)
            if error != None:
                return jsonify({"error": error})
            forecasted_earning, error = forecast(total_earned, int(months), 1, 1, 1)
            recent_month = dates[-1].get('month')
            recent_year = dates[-1].get('year')

            res_spent = []
            res_earned = []
            for i in range(len(forecasted_spending)):
                month = (recent_month + i+1) % 12
                year = recent_year + (recent_month + i+1) // 12
                res_spent.append({"Month": month, "Year": year, "Amount": forecasted_spending[i]})
                res_earned.append({"Month": month, "Year": year, "Amount": forecasted_earning[i]})
            sub_response = {}
            sub_response['item_name'] = item.get("ItemName")
            sub_response["total_spending"] = item.get("TotalSpent")
            sub_response["total_earning"] = item.get("TotalEarned")
            sub_response["forecasted_spending"] = res_spent
            sub_response["forecasted_earning"] = res_earned

            sub_response["net_cash_flow"] = float(sum(forecasted_earning) - sum(forecasted_spending))
            sub_response["recommended_budget"] = float(mean_value(total_spent + forecasted_spending))
            response.append(sub_response)
        return jsonify(response)

api.add_resource(Forecast, '/forecast')

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)