from flask import Flask, request, jsonify
from flask_restful import Api, Resource
from ARIMA import forecast
import numpy as np
from flask_cors import CORS, cross_origin

app = Flask(__name__)

api = Api(app)
CORS(app, supports_credentials=True)

class Forecast(Resource):
    @cross_origin()
    def post(self):
        jsonData = request.get_json()
        print(jsonData)
        # get the number of months to forecast from url
        months = request.args.get('months')
        if not jsonData:
            return jsonify({"message": "cannot find JSON data provided"}, 400)

        transactions = []
        dates = []
        # append all transaction amount to transactions array
        for transaction in jsonData:
            month = transaction.get('Month')
            year = transaction.get('Year')
            total_amount = transaction.get('TotalAmount')
            transactions.append(total_amount)
            dates.append({'month': month, 'year': year})
        
        # forecast transactions
        forecasted_transactions = forecast(transactions, int(months))

        if np.array_equal(np.array(forecasted_transactions), np.array([-1])):
            error = jsonify({
                'message': '''There are not enough transactions to make an accurate forecast, 
                please provide at least 5 months of transactions'''})
            error.status_code = 400
            return error

        # generate month and year for forecasted transactions
        latest = max(dates, key = lambda x: (x['year'], x['month']))
        latest_month = latest['month']
        latest_year = latest['year']

        res = []
        # append all forecasted transactions and their date onto res
        for i in range(len(forecasted_transactions)):
            forecast_month = (latest_month + i+1) % 12
            forecast_year = latest_year + (latest_month + i+1) // 12

            if forecast_month == 0:
                forecast_month = 12
                forecast_year -= 1

            res.append({
                'Month': forecast_month,
                'Year': forecast_year,
                'TotalAmount': forecasted_transactions[i],
            })

        #combine the original transactions with the forecast
        combined_transactions = np.concatenate((transactions, forecasted_transactions))

        #approximating the original and forecasted transactions as a polynomial
        x = np.arange(len(combined_transactions))
        coefficients = np.polyfit(x, combined_transactions, deg=7)
        fitted_poly = np.poly1d(coefficients)

        #integrate the polynomial
        integral = np.polyint(fitted_poly)
        first = x[0]
        last = x[-1]

        #calculate the mean value of the transactions and set it as recommended budget
        mean_value = (integral(last) - integral(first)) / (last - first)
        recommended = mean_value

        return jsonify({
            'total_transactions': jsonData,
            'forecast': res,
            'recommended_budget': recommended
        })

api.add_resource(Forecast, '/forecast')

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)