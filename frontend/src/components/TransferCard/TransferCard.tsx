import type { Transfer } from "../../api/types";
import { formatDateTime, formatMoney } from "../../utils/format";
import styles from "./TransferCard.module.css";

interface TransferCardProps {
    transfer: Transfer;
}

export function TransferCard({transfer}: TransferCardProps) {
    return (
        <div className={styles.card}>
            <div>
                <span>№</span>
                <span>{transfer.account_from_id}</span>
                <span>➜</span>
                <span>№</span>
                <span>{transfer.account_to_id}</span>
            </div>

            <div>
                <span>Сумма перевода: </span>
                <span>{formatMoney(transfer.amount)}</span>
            </div>

            <div>
                <span>Дата: </span>
                <span>{formatDateTime(transfer.created_at)}</span>
            </div>
        </div>
    )
}