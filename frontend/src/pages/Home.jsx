import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Bar, Line } from "react-chartjs-2";
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
    let redirect = useNavigate();
    const token = localStorage.getItem('token');
    const username = localStorage.getItem('username');

    const [transactions, setTransactions] = useState([]);
    const [chartToggle, setChartToggle] = useState("bar");
    const [chartData, setChartData] = useState({ labels: [], datasets: [] });
    const [options, setOptions] = useState()
    // options for displaying the line chart: x-axis: Transaction date, y-axis: amount
    const optionLine = {
        responsive: true,
        plugins: {
            legend: { position: 'top' },
            tooltip: {
                callbacks: {
                    label: (tooltipItem) => {
                        return `${tooltipItem.dataset.label}: ${tooltipItem.raw}`;
                    }
                }
            }
        },
        scales: {
            x: { stacked: true, title: { display: true, text: 'Transaction date' } },
            y: { stacked: true, title: { display: true, text: 'Transaction amount' } }
        }
    };

    // options for stacked bar chart: x-axis: budgets, y-axis: total transaction amount
    const optionBar = {
        responsive: true,
        plugins: {
            legend: { position: 'top' },
            tooltip: {
                callbacks: {
                    label: (tooltipItem) => {
                        return `${tooltipItem.dataset.label}: ${tooltipItem.raw}`;
                    }
                }
            }
        },
        scales: {
            x: { stacked: true, title: { display: true, text: 'Budgets' } },
            y: { stacked: true, title: { display: true, text: 'Total transaction amount' } }
        }
    };

    useEffect(() => {
        if (token === null) {
            redirect('/login');
        }

        // make sure get transactions is asynchronous
        const getTransactions = async () => {
            let date = new Date();
            // get request to retrieve all transactions of the current month
            await axios.get(
                `http://localhost:8080/main/transactions/${date.getFullYear()}/${date.getMonth() + 1}`, 
                { headers: { Authorization: `Bearer ${token}` } }
            ).then(response => {
                // check if response data is null or is an array
                if (!response.data || !Array.isArray(response.data)) {
                    setTransactions([]);
                    return;
                }
                setTransactions(response.data);
                // initially group transactions as stacked bar chart
                groupTransactionsIntoStackedBarChart(response.data);
            }).catch(error => { // return error message if there is any
                alert(error.response?.data || error.message); 
            })

        };
        getTransactions();
    }, [redirect, token]);

    const groupTransactionsIntoLineChart = (data) => {
        // return nothing is there is no data present
        if (!data || data.length === 0) return

        // extract unique transaction dates from each transactions
        // and sort then in date order
        const labels = Array.from(new Set(
            data.flatMap(d => [...d.transactions].map(t => 
                new Date(t.date).toISOString().split('T')[0]
            ).sort((a, b) => new Date(a) - new Date(b)))
        ))

        const datasets = data.flatMap(d => {
            const colour = getRandomColour()
            return [
                {
                    // make label for inflow transactions
                    label: `${d.item.item_name} - inflow`,
                    data: labels.map(label => {
                    // match the transaction to the correct x value
                    const entry = d.transactions.find(t => 
                        new Date(t.date).toISOString().split('T')[0] === label && t.type === "inflow")
                    return entry ? entry.amount : null
                    }),
                    borderColor: colour,
                    backgroundColor: "rgba(0,0,0,0)",
                    tension: 0.3,
                    // line is dotted
                    borderDash: [5, 5]
                },
                {
                // make label for outflow transactions
                label: `${d.item.item_name} - outflow`,
                data: labels.map(label => {
                    // match the transaction to the correct x value
                    const entry = d.transactions.find(t => 
                        new Date(t.date).toISOString().split('T')[0] === label && t.type === "outflow")
                    return entry ? entry.amount : null
                }),
                borderColor: colour,
                backgroundColor: "rgba(0,0,0,0)",
                tension: 0.3,
                // line is solid
                borderDash: []
                }
            ]
        })

        // replace the options with the options for a line chart
        setOptions(optionLine)
        setChartData({labels, datasets})
    }

    const groupTransactionsIntoStackedBarChart = (data) => {
        // return nothing is there is no data present
        if (!data || data.length === 0) return

        let budget_data_map = {}

        data.forEach(t => {
            // return nothing if transaction data is empty
            if (!t || !t.item || !t.transactions) return;
            const budget_name = t.item.budget_name;
            const item_name = t.item.item_name;

            // create a key value pair for budget name 
            // if it doesn't exist in the hashmap
            if (!budget_data_map[budget_name]) {
                budget_data_map[budget_name] = {};
            }
            // create a key value pair for item name 
            // inside its corresponding budget name
            // if it doesn't exist in the hashmap
            if (!budget_data_map[budget_name][item_name]) {
                budget_data_map[budget_name][item_name] = { inflow: 0, outflow: 0 };
            }
            t.transactions.forEach(transaction => {
            if (!transaction) return;
            // checks if the current transaction is an inflow 
            // or outflow
            switch (transaction.type) {
                case "outflow":
                    // increment the amount to inflow
                    budget_data_map[budget_name][item_name].outflow += transaction.amount;
                    break
                case "inflow":
                    // increment the amount to outflow
                budget_data_map[budget_name][item_name].inflow += transaction.amount;
                    break
                default:
                    break
                }            
            })
        });

        // extract the keys (budget names) in the hashmap
        const labels = Object.keys(budget_data_map);
        const budget_items = new Set();

        // add the corresponding items to the budgets
        Object.values(budget_data_map).forEach(items => {
            Object.keys(items).forEach(item => budget_items.add(item));
        });

        // arrange the data into a the dataset
        // and set their colours
        const datasets = [
            ...Array.from(budget_items).map((item) => ({
                label: `${item} (inflow)`,
                data: labels.map(budget => budget_data_map[budget]?.[item]?.inflow || 0),
                backgroundColor: getRandomColour(),
            })),
            ...Array.from(budget_items).map((item) => ({
                label: `${item} (outflow)`,
                data: labels.map(budget => budget_data_map[budget]?.[item]?.outflow || 0),
                backgroundColor: getRandomColour(),
            }))
        ];

        // replace the options with options for a stacked bar chart
        setOptions(optionBar)
        setChartData({ labels, datasets });  // Update chartData state
    };

    // toggles between line and stacked bar depending 
    // on the user's input
    useEffect(() => {
        if (chartToggle === "line") {
            groupTransactionsIntoLineChart(transactions)
        } else {
            groupTransactionsIntoStackedBarChart(transactions)
        }
    }, [transactions, chartToggle])

    const handleToggleChart = () => {
        setChartToggle(prevToggle => {
            const newToggle = prevToggle === "bar" ? "line" : "bar";
            return newToggle;
        });
    };    

    const getRandomColour = () => {
        const letters = '0123456789ABCDEF';
        let colour = '#';
        for (let i = 0; i < 6; i++) {
            colour += letters[Math.floor(Math.random() * 16)];
        }
        return colour;
    };

    return (
        <div className="p-20">
            <div className="flex justify-evenly">
                <h1>Welcome {username}</h1>
            </div>
            <h1>Summary of transactions this month:</h1>
            {chartData.labels.length > 0 ? (
                chartToggle === "bar" ? <Bar className="p-20" data={chartData} options={options} />
                : <Line className="p-20" data={chartData} options={options} />
            ) : (
                <p>You have made no transactions</p>
            )}
            <button onClick={handleToggleChart}>
                Toggle Chart
            </button>
            <div className="flex justify-evenly">
                <button onClick={(e) => { e.preventDefault(); redirect('/budgets') }}>View budgets</button>
                <button onClick={(e) => { e.preventDefault(); redirect('/transactions') }}>View transactions</button>
                <button onClick={(e) => { e.preventDefault(); redirect('/forecast') }}>Forecast transactions</button>
            </div>
        </div>
    );
};

export default Home;