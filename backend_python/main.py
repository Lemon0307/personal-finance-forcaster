from flask import Flask, request, jsonify
from flask_restful import Api, Resource
from ARIMA import forecast
import numpy as np

app = Flask(__name__)
api = Api(app)

class Forecast(Resource):
    def post(self): 
        jsonData = request.get_json()
        months = request.args.get('months')
        if not jsonData:
            return jsonify({"message": "cannot find JSON data provided"}, 400)

        transactions = []
        dates = []
        for transaction in jsonData:
            month = transaction.get('Month')
            year = transaction.get('Year')
            total_amount = transaction.get('TotalAmount')
            transactions.append(total_amount)
            dates.append({'month': month, 'year': year})
        
        forecasted_transactions = forecast(transactions, int(months))

        latest = max(dates, key = lambda x: (x['year'], x['month']))
        latest_month = latest['month']
        latest_year = latest['year']

        res = []
        for i in range(len(forecasted_transactions)):
            forecast_month = (latest_month + i+1) % 12
            forecast_year = latest_year + (latest_month + i+1) // 12

            if forecast_month == 0:
                forecast_month = 12
                forecast_year -= 1

            res.append({
                'month': forecast_month,
                'year': forecast_year,
                'forecasted_transaction': forecasted_transactions[i]
            })

        return jsonify(res)

class RecommendBudget(Resource):
    def post(self): 
        jsonData = request.get_json()
        if not jsonData:
            return jsonify({"message": "cannot find JSON data provided"}, 400)

        transactions = []
        for transaction in jsonData:
            total_amount = transaction.get('TotalAmount')
            transactions.append(total_amount)
        
        recommendation = forecast(transactions, 3)
        avg = np.mean(recommendation)

        return jsonify({
            'recommended_budget_amount': avg
        })

api.add_resource(Forecast, '/forecast/')
api.add_resource(RecommendBudget, '/recommend_budget/')

if __name__ == '__main__':
    app.run(debug=True, host='localhost', port=5000)