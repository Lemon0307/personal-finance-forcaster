import React, { useState, useRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Bar } from "react-chartjs-2";
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend
  } from 'chart.js';
import axios from "axios";

ChartJS.register(
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend
  );

const Home = () => {
    let redirect = useNavigate()
    const token = localStorage.getItem('token')
    const username = localStorage.getItem('username')
    const [data, setData] = useState(null)

    // transactions charts options
    const options = {
        responsive: true,
        plugins: {
            legend: {
                position: 'top',
            },
            tooltip: {
                callbacks: {
                    label: function(tooltipItem) {
                        return `Transactions: ${tooltipItem.raw}`;
                    }
                }
            }
        },
        scales: {
            x: {
                stacked: true,
                title: {
                    display: true,
                    text: 'Budgets'
                }
            },
            y: {
                stacked: true,
                title: {
                    display: true,
                    text: 'Total transaction amount'
                }
            }
        }
    };

    useEffect(() => {
        // logout user if token is missing
        if (token === null) {
            redirect('/login')
        }
        
        const getTransactions = async () => {
            let date = new Date()
            // get transactions of the current month
            await axios.get(`http://localhost:8080/main/transactions/${date.getFullYear()}/${date.getMonth() + 1}`, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            })
            .then(response => {
                if (response.data && (Array.isArray(response.data) || response.data.length > 0)) { // check if there are transactions in the response
                    // reduce the response into a single object
                    const budgetData = response.data.reduce((acc, budget) => {
                        const budget_name = budget.item.budget_name;
                        const item_name = budget.item.item_name;
                        // sums all transactions in the transaction array, if there isn't any in the array then return zero
                        const total_amount = budget.transactions?.reduce((sum, transaction) => sum + transaction.amount, 0) ?? 0;

                        acc[budget_name] ??= {}
                        // store total transaction amount in each item
                        acc[budget_name][item_name] = total_amount;
        
                        return acc;
                    }, {});

                    // get all budgets
                    const all_budgets = Object.keys(budgetData);
                    // get all unique items
                    const all_items = new Set();
                    all_budgets.forEach(budget => {
                        Object.keys(budgetData[budget]).forEach(item => {
                            all_items.add(item);
                        });
                    });
                    // convert to array
                    const item_names = Array.from(all_items);

                    // perpare chart data
                    const labels = all_budgets;  // labels x-axis
                    const datasets = item_names.map(item_name => { // labels y-axis
                        return {
                            label: item_name,
                            data: all_budgets.map(budget => budgetData[budget][item_name]),
                            backgroundColor: getRandomColour()
                        };
                    });

                    // group x and y axis to form the chart
                    const chartData = {
                        labels: labels,
                        datasets: datasets,
                    };
                    setData(chartData)
                }
            }).catch(error => { // return error message
                alert(error.response?.data || error.message);
            })
        }
        getTransactions();
    }, [redirect, token]);

    const getRandomColour = () => {
        const letters = '0123456789ABCDEF';
        let colour = '#';
        for (let i = 0; i < 6; i++) {
          colour += letters[Math.floor(Math.random() * 16)];
        }
        return colour;
    }

    return (
        <div className="p-20">
            <div className="flex justify-evenly">
                <h1>Welcome {username}</h1>
            </div>
            <h1>Summary of transactions this month:</h1>{
            data ? (
                <Bar className="p-20" data={data} options={options} />
            ) : (
                <p>You have made no transactions</p>
            )}
            <div className="flex justify-evenly">
                <button onClick={(e) => {e.preventDefault(); redirect('/budgets')}}>View budgets</button>
                <button onClick={(e) => {e.preventDefault(); redirect('/transactions')}}>View transactions</button>
                <button onClick={(e) => {e.preventDefault(); redirect('/forecast')}}>Forecast transactions</button>
            </div>
        </div>
    )
    
}

export default Home;