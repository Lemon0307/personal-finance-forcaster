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


    for i in range(1, len(data)):
        # linear combination of phi and theta with each data and its error
        prediction[i] = np.dot([a, b], [data[i-1], error[i-1]])
        error[i] = data[i] - prediction[i]

    last = data[-1]
    last_error = error[-1]

    for i in range(months):
        forecast[i] = np.dot([a, b], [last, last_error])
        last += forecast[i]
        last_error = forecast[i] - a*last

    print(forecast)
    return np.array(forecast)

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