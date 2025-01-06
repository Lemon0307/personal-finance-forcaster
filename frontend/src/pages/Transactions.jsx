import { useState, useEffect } from "react";
import { useParams } from 'react-router-dom';

const Transactions = () => {
    useEffect(async () => {
        
    })
    const {item_name} = useParams()

    return (
        <div className="p-20 flex">
            <table>
                <tr>
                    <th>Name</th>
                    <th>Type</th>
                    <th>Amount</th>
                    <th>Date</th>
                </tr>
            </table>
        </div>
    )
}

export default Transactions;