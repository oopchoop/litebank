export interface Account {
    id: number;
    account_type_id: number;
    balance: number;
    is_system: boolean;
    created_at: string;
}

export interface Transfer {
    id: number;
    amount: number;
    account_from_id: number;
    account_to_id: number;
    created_at: string;
}

export interface TransferIncoming extends Transfer {
    is_system: boolean;
}

export interface TransferList {
    transfers_sent: Transfer[];
    transfers_incoming: TransferIncoming[];
}