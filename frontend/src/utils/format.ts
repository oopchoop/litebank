export const formatMoney = (sum: number) => {
      return (sum / 100).toLocaleString('ru-RU', {
            style: 'currency',
            currency: 'RUB',
            minimumFractionDigits: 2,
      });
}

export const formatDate = (dateString: string) => {
    if (!dateString) {
        return '';
    }

    const date = new Date(dateString);

    return date.toLocaleDateString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
    });
}

export const formatDateTime = (dateString: string) => {
    if (!dateString) {
        return '';
    }

    const date = new Date(dateString);

    return date.toLocaleString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    })
}