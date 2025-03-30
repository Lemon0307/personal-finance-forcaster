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

    const options = {
        responsive: true,
        plugins: {
            legend: { position: 'top' },
            tooltip: {
                callbacks: {
                    label: function(tooltipItem) {
                        return `Transactions: ${tooltipItem.raw}`;
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

        const getTransactions = async () => {
            let date = new Date();
            try {
                const response = await axios.get(
                    `http://localhost:8080/main/transactions/${date.getFullYear()}/${date.getMonth() + 1}`, 
                    { headers: { Authorization: `Bearer ${token}` } }
                );
                setTransactions(response.data);
                console.log(response.data)
                groupTransactionsIntoStackedBarChart(response.data);  // Ensure data updates correctly
            } catch (error) {
                alert(error.response?.data || error.message);
            }
        };

        getTransactions();
    }, [redirect, token]);

    const groupTransactionsIntoLineChart = (data) => {
        if (!data || data.length === 0) return

        const labels = Array.from(new Set(
            data.flatMap(d => [...d.transactions].map(t => 
                new Date(t.date).toISOString().split('T')[0]
            ).sort((a, b) => new Date(a) - new Date(b)))
        ))

        const datasets = data.flatMap(d => {
            const colour = getRandomColour()
            return [
                {
                 label: `${d.item.item_name} - inflow`,
                 data: labels.map(label => {
                    const entry = d.transactions.find(t => 
                        new Date(t.date).toISOString().split('T')[0] === label && t.type === "inflow")
                    return entry ? entry.amount : null
                 }),
                 borderColor: colour,
                 backgroundColor: "rgba(0,0,0,0)",
                 tension: 0.3,
                 borderDash: [5, 5]
                },
                {
                label: `${d.item.item_name} - outflow`,
                data: labels.map(label => {
                    const entry = d.transactions.find(t => 
                        new Date(t.date).toISOString().split('T')[0] === label && t.type === "outflow")
                    return entry ? entry.amount : null
                }),
                borderColor: colour,
                backgroundColor: "rgba(0,0,0,0)",
                tension: 0.3,
                borderDash: []
                }
            ]
        })

        setChartData({labels, datasets})
    }

    const groupTransactionsIntoStackedBarChart = (data) => {
        if (!data || data.length === 0) return

        let budget_data_map = {}

        data.forEach(t => {
            const budget_name = t.item.budget_name;
            const item_name = t.item.item_name;

            if (!budget_data_map[budget_name]) {
                budget_data_map[budget_name] = {};
            }
            if (!budget_data_map[budget_name][item_name]) {
                budget_data_map[budget_name][item_name] = { inflow: 0, outflow: 0 };
            }
            t.transactions.forEach(transaction => {
            if (transaction.type === "outflow") {
                budget_data_map[budget_name][item_name].outflow += transaction.amount;
            } else if (transaction.type === "inflow") {
                budget_data_map[budget_name][item_name].inflow += transaction.amount;
            }                
            })
        });

        const labels = Object.keys(budget_data_map);
        const budget_items = new Set();
        Object.values(budget_data_map).forEach(items => {
            Object.keys(items).forEach(item => budget_items.add(item));
        });

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

        console.log(budget_data_map)
        setChartData({ labels, datasets });  // Update chartData state
    };

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