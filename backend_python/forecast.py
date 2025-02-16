import pandas as pd
import numpy as np
from scipy.stats import boxcox
from statsmodels.tsa.arima.model import ARIMA
from scipy.special import inv_boxcox

def forecast(data, months):
    # turns the data into a series
    data = pd.Series(data)
    # applies the boxcox transform to stationarise data
    stationary_data, lam = boxcox(data)
    # difference the data for further stationarity
    differenced_data = pd.Series(stationary_data).diff().dropna()
    # build and apply the ARIMA model for forecasting
    model = ARIMA(stationary_data, order=(12, 1, 12)).fit()
    boxcox_forecasts = model.forecast(steps=months)
    #apply inverse boxcox to revert the data back
    forecasted_values = inv_boxcox(boxcox_forecasts, lam)

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