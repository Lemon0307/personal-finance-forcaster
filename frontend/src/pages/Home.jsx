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
    const [transactions, setTransactions] = useState()
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
                setTransactions(response.data)
            }).catch(error => { // return error message
                alert(error.response?.data || error.message);
            })
        }

        getTransactions();
    }, [redirect, token]);

    const groupTransactionsIntoLineChart = (data) => {

    }

    const groupTransactionsIntoStackedBarChart = (data) => {
    }

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