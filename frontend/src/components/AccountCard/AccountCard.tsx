import type { Account } from "../../api/types";
import { formatDate, formatMoney } from "../../utils/format";
import { Button } from "../Button/Button";
import styles from "./AccountCard.module.css";

interface AccountCardProps {
    account: Account;
    onDepositClick: (accountId: number) => void;
    onTransferClick: (accountId: number) => void;
}

export function AccountCard({account, onDepositClick, onTransferClick}: AccountCardProps) {
    return (
        <div className={styles.card}>
            <div>
                <span>Личный счет: </span>
                <span>{account.id}</span>
            </div>

            <div>
                <span>Баланс: </span>
                <span>{formatMoney(account.balance)}</span>
            </div>

            <div>
                <span>Создан: </span>
                <span>{formatDate(account.created_at)}</span>
            </div>

            <div className={styles.footer}>
                <Button onClick={() => onDepositClick(account.id)}>Пополнить счет</Button>
                <Button onClick={() => onTransferClick(account.id)}>Перевести</Button>
            </div>
        </div>
    )
}