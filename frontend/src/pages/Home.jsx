import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Line, Bar } from "react-chartjs-2";
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
    
    useEffect(() => {
        if (localStorage.getItem('token') !== null) {
            const getTransactions = async () => {
                let date = new Date()
                const token = localStorage.getItem('token')
                try {
                    const response = await axios.get(`http://localhost:8080/main/transactions/${date.getFullYear()}/${date.getMonth() + 1}`, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });

                    console.log(response.data)

                    if (!Array.isArray(response.data) || response.data.length === 0) {
                        setHasTransactions(false);
                        setData(null);
                        return;
                    }
        
                    const groupedData = response.data.reduce((acc, budget) => {
                        const budgetName = budget.budget_item.budget_name;
                        const itemName = budget.budget_item.item_name;
                        const transactionAmount = budget.transactions ? budget.transactions.reduce((sum, transaction) => sum + transaction.amount, 0) : 0;

                        if (!acc[budgetName]) {
                            acc[budgetName] = {};
                        }
        
                        // Store transaction counts for each item under a budget
                        acc[budgetName][itemName] = transactionAmount;
        
                        return acc;
                    }, {});

                // Get unique budget names and item names
                const budgetNames = Object.keys(groupedData);
                const allItemNames = new Set();
                budgetNames.forEach(budget => {
                    Object.keys(groupedData[budget]).forEach(item => {
                        allItemNames.add(item);
                    });
                });

                const itemNames = Array.from(allItemNames);  // Convert set to array

                // Prepare data for chart
                const labels = budgetNames;  // X-axis will have budget names
                const datasets = itemNames.map(itemName => {
                    return {
                        label: itemName,
                        data: budgetNames.map(budget => groupedData[budget][itemName] || 0), // For each budget, get the count for this item
                        backgroundColor: getRandomColour(),
                        borderColor: getRandomColour(),
                        borderWidth: 1,
                    };
                });

                const chartData = {
                    labels: labels,  // X-axis: Budget Names
                    datasets: datasets,  // Y-axis: Transaction counts for each item
                };

                setHasTransactions(true);
                setData(chartData);  // Set chart data for rendering
        
                } catch (error) {
                    alert(error.response?.data || error.message);
                }
            }
        
            getTransactions();
        } else {
            redirect('/login')
        }
    }, [redirect]);

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

    const username = localStorage.getItem('username')
    const [choice, setChoice] = useState()
    const [data, setData] = useState(null)
    const [hasTransactions, setHasTransactions] = useState(true);

    const getRandomColour = () => {
        const letters = '0123456789ABCDEF';
        let colour = '#';
        for (let i = 0; i < 6; i++) {
          colour += letters[Math.floor(Math.random() * 16)];
        }
        return colour;
    }

    const choiceHandler = (e) => {
        setChoice(e.target.value)
    }

    return (
        <div className="p-20">
            <div className="flex justify-evenly">
                <h1>Welcome {username}</h1>
            </div>
            <h1>Summary of transactions {choice}:</h1>
            {hasTransactions ? (
                data ? (
                    <Bar className="p-20" data={data} options={options} />
                ) : (
                    <p>Loading transactions...</p>
                )
            ) : (
            <p>You have made no transactions {choice}</p>
            )}
            <select value={choice} onChange={choiceHandler}>
                <option value="this week">This week</option>
                <option value="this month">This month</option>
            </select>
            <div className="flex justify-evenly">
                <button onClick={(e) => {e.preventDefault(); redirect('/budgets')}}>View budgets</button>
                <button onClick={(e) => {e.preventDefault(); redirect('/transactions')}}>View transactions</button>
                <button onClick={(e) => {e.preventDefault(); redirect('/forecast')}}>Forecast transactions</button>
            </div>
        </div>
    )
    
}

export default Home;