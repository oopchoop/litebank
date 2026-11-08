import type { Transfer, TransferIncoming } from "../../api/types";
import { formatDateTime, formatMoney } from "../../utils/format";
import styles from "./TransferCard.module.css";

interface TransferCardProps {
    transfer: Transfer | TransferIncoming;
}

export function TransferCard({transfer}: TransferCardProps) {
    const isSystem = "is_system" in transfer && transfer.is_system;

    const sumText = isSystem ? 'Сумма пополнения: ' : 'Сумма перевода: ';

    return (
        <div className={styles.card}>
            <div>
                {isSystem ? (
                    <span>Пополнение баланса</span>
                ) : (
                    <>
                        <span>№</span>
                        <span>{transfer.account_from_id}</span>
                        <span>➜</span>
                        <span>№</span>
                        <span>{transfer.account_to_id}</span>
                    </>                
                )}
            </div>

            <div>
                <span>{sumText}</span>
                <span>{formatMoney(transfer.amount)}</span>
            </div>

            <div>
                <span>Дата: </span>
                <span>{formatDateTime(transfer.created_at)}</span>
            </div>
        </div>
    )
}