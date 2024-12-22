import numpy as np

def forecast(data, months):
    data = np.array(data)
    
    # calculates the n-th discrete difference along the given axis.
    difference = np.diff(data)

    #finds a suitable value for phi
    a = estimate_first_ar(difference)

    #finds a suitable value for theta
    b = estimate_first_ma(difference, a)

    # the main ARIMA model
    forecast = ARIMA(difference, a, b, months)

    # returns the last value of data and forecast as an array
    return forecast + data[-1]

def ARIMA(data, a, b, months):
    
    prediction = np.zeros(len(data))

    forecast = np.zeros(months)

    error = np.zeros(len(data))

    lagged_matrix = np.column_stack((data[:-1], error[:-1]))
    prediction[1:] = lagged_matrix @ np.array([a, b])
    error[1:] = data[1:] - prediction[1:]

    last = data[-1]
    last_error = error[-1]

    lagged_forecast = np.zeros((months, 2))
    lagged_forecast[0] = [data[-1], error[-1]]
    for i in range(1, months):
        lagged_forecast[i] = [lagged_forecast[i-1] @ np.array([a, b]),
        lagged_forecast[i-1, 0]]
    forecast[:] = lagged_forecast[:, 0]

    return forecast

def estimate_first_ar(data):
    lag_1 = data[:-1]
    y_t = data[1:]
    phi = np.corrcoef(lag_1, y_t)[0, 1]
    return phi

def estimate_first_ma(data, a):
    predicted_ar = np.roll(data, 1) * a
    error = data - predicted_ar
    error = error[1:]
    
    theta_1 = np.corrcoef(error[:-1], error[1:])[0, 1]
    return theta_1