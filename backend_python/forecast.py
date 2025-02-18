import pandas as pd
import numpy as np
from scipy.stats import boxcox
from statsmodels.tsa.arima.model import ARIMA
from scipy.special import inv_boxcox

def forecast(transactions, months, p, d, q):
    # turns the transactions into a series
    transactions = pd.Series(transactions)
    # applies the boxcox transform to stationarise transactions
    stationary_data, lam = boxcox(transactions)
    # difference the transactions for further stationarity
    differenced_data = pd.Series(stationary_data).diff().dropna()
    # build and apply the ARIMA model for forecasting
    model = ARIMA(differenced_data, order=(p, d, q)).fit()
    differenced_forecasts = model.forecast(steps=months)
    #reverse differencing
    last = stationary_data[-1]
    forecasted_stationary = np.r_[last, differenced_forecasts].cumsum()[1:]
    #apply inverse boxcox to revert the transactions back
    forecasted_values = inv_boxcox(forecasted_stationary, lam)

    return forecasted_values

def mean_value(combined_transactions):
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

    return mean_value

def knapsack(values, weights, capacity):
    n = len(values)
    dp = [0] * (capacity + 1)

    for i in range(n):
        for w in range(capacity, weights[i] - 1, -1):
            dp[w] = max(dp[w], dp[w - weights[i]] + values[i])

    # Backtrack to find which items are selected
    w = capacity
    selected_items = [False] * n
    for i in range(n - 1, -1, -1):
        if w >= weights[i] and dp[w] == dp[w - weights[i]] + values[i]:
            selected_items[i] = True
            w -= weights[i]

    return selected_items