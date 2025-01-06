import { useState, useEffect } from "react";
import { useParams } from 'react-router-dom';

const Transactions = () => {
    // useEffect(async () => {
        
    // })
    const {item_name} = useParams()
    const [sort, setSort] = useState()

    const handleSort = (e) => {
        setSort(e.target.value)
        switch (sort) {

        }
    }

    return (
        <div className="p-20 flex">
            <div className="px-10">
                <h1>Sort by:</h1>
                <select value={sort} onChange={handleSort}>
                    <option value="name">Name</option>
                    <option value="type">Type</option>
                    <option value="amount">Amount</option>
                    <option value="date">Date</option>
                </select>
            </div>
            <div className="grid place-content-center">
                <h1>My Transactions</h1>
                <table className="flex justify-normal">
                    <tr>
                        <th className="px-5">Name</th>
                        <th className="px-5">Type</th>
                        <th className="px-5">Amount</th>
                        <th className="px-5">Date</th>
                    </tr>
                </table>                
            </div>

        </div>
    )
}

export default Transactions;