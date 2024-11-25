from datetime import datetime, timezone, timedelta

def get_execution_date(date_input=None):
    """
    Determina a data de execução com base no input:
    - Se None: usa data atual + 5 minutos
    - Se string: converte para datetime UTC
    - Se 'morning': próximo dia às 9h UTC
    """
    current_date = datetime.now(timezone.utc)

    if date_input is None:
        # Caso padrão: 5 minutos a partir de agora
        return current_date + timedelta(minutes=5)

    if date_input == 'morning':
    # Próximo dia às 9h UTC
        return (current_date + timedelta(days=1)).replace(hour=9, minute=0, second=0, microsecond=0)
    
    if isinstance(date_input, str):
        # Converte string no formato ISO para datetime UTC
        return datetime.strptime(date_input, '%Y-%m-%dT%H:%MZ').replace(tzinfo=timezone.utc)
    
    # Se for passado um datetime diretamente
    if isinstance(date_input, datetime):
        # Garante que está em UTC
        return date_input.astimezone(timezone.utc)
    
    raise ValueError("Input inválido para data de execução")