import numpy as np

def forecast(data, months):
    data = np.array(data)
    
    # calculates the difference between data points
    difference = np.diff(data)
    
    #return a number (-1) if there are nan of infinite numbers in difference
    if np.any(np.isnan(difference)) or np.any(np.isinf(difference)):
        return [-1]

    #finds a suitable value for phi
    a = estimate_first_ar(difference)

    #finds a suitable value for theta
    b = estimate_first_ma(difference, a)

    # the main ARIMA model
    forecast = ARIMA(difference, a, b, months)

    #return a number (-1) if there are nan of infinite numbers in forecast
    if np.any(np.isnan(forecast)) or np.any(np.isinf(forecast)):
        return [-1]

    # returns the last value of data and forecast as an array
    return forecast + data[-1]

def ARIMA(data, a, b, months):
    
    prediction = np.zeros(len(data))
    forecast = np.zeros(months)
    error = np.zeros(len(data))

    # create a matrix for current data and errors
    lagged_matrix = np.column_stack((data[:-1], error[:-1]))

    # predict for the given data
    prediction[1:] = lagged_matrix @ np.array([a, b])

    # calculate the error for each data point
    error[1:] = data[1:] - prediction[1:]

    # store last value and error for forecasting
    last = data[-1]
    last_error = error[-1]

    # create another matrix for forecasting
    lagged_forecast = np.zeros((months, 2))
    lagged_forecast[0] = [last, last_error]

    # forecast future values using the params a and b
    for i in range(1, months):
        lagged_forecast[i] = [lagged_forecast[i-1] @ np.array([a, b]),
        lagged_forecast[i-1, 0]]

    # store forecasted values from matrix into forecast array
    forecast[:] = lagged_forecast[:, 0]

    return forecast

# def estimate_first_ar(difference):
#     previous_values = difference[:-1]
#     current_values = difference[1:]

#     # calculates the r value for the data sets
#     # which is the value of phi
#     phi = np.corrcoef(previous_values, current_values)[0, 1]
#     return phi

# def estimate_first_ma(difference, a):
#     # predict AR component
#     predicted_ar = np.roll(difference, 1) * a

#     # calculate errors
#     error = difference - predicted_ar
#     error = error[1:] # first value doesn't matter so we drop it
        
#     # estimate MA coefficient using PCC and calculate the r value
#     theta = np.corrcoef(error[:-1], error[1:])[0, 1]
#     return theta

def estimate_first_ar(difference):
    """ Estimate AR(1) coefficient using Yule-Walker but prevent extreme values. """
    n = len(difference)
    if n < 2:
        return 0  # Avoid errors if not enough data

    y = difference[1:]
    X = difference[:-1]

    # Solve for phi using least squares
    phi = np.dot(X, y) / np.dot(X, X)

    # Prevent extreme values that cause sharp drops
    phi = max(min(phi, 0.9), -0.9)

    return phi

def estimate_first_ma(difference, a):
    """ Estimate MA(1) coefficient using error minimization. """
    predicted_ar = np.roll(difference, 1) * a
    error = difference - predicted_ar
    error = error[1:]  # Drop the first value

    # Solve for theta using least squares
    theta = np.dot(error[:-1], error[1:]) / np.dot(error[:-1], error[:-1])

    # Prevent extreme values that amplify trends
    theta = max(min(theta, 0.9), -0.9)

    return theta