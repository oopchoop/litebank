import api from "../../api/client";
import { useState, useEffect } from "react";
import type { Transfer, TransferIncoming, TransferList } from "../../api/types";
import { TransferCard } from "../../components/TransferCard/TransferCard";
import styles from "./Histroy.module.css";


export function History() {
    const [transfersIncoming, setTransfersIncoming] = useState<TransferIncoming[]>([]);
    const [transfersSent, setTransfersSent] = useState<Transfer[]>([]);
    const [activeTab, setActiveTab] = useState<"incoming" | "sent">("incoming");

    const fetchTransfers = async () => {
        try {
            const response = await api.get<TransferList>('/transfers')

            const responseData = response.data;

            setTransfersIncoming(responseData.transfers_incoming)
            setTransfersSent(responseData.transfers_sent)
        } catch (err: any) {
            console.log(err);
        }
    }

    useEffect(() => {
        fetchTransfers()
    }, [])

    const currentTransfers = activeTab === "incoming" ? transfersIncoming : transfersSent;
    const currentError = activeTab === "incoming" ? "You don`t have incoming transfers" : "You don`t have sent transfers";

    return (
        <div className={styles.container}>
            <div className={styles.tabs}>
                <button className={activeTab === "incoming" ? styles.active : ""}
                        onClick={() => setActiveTab("incoming")}>
                    Входящие
                </button>
                <button className={activeTab === "sent" ? styles.active : ""}
                        onClick={() => setActiveTab("sent")}>
                    Исходящие
                </button>
            </div>
            <div className={styles.transfers}>
                {currentTransfers.length === 0 ? (
                    <div>{currentError}</div>   
                ) : (
                    currentTransfers.map(item => (
                        <TransferCard key={item.id}
                                    transfer={item} />
                    ))
                )}
            </div>  
        </div>
    )
}