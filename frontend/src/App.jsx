import {useState} from 'react';
import {DecryptFile, EncryptFile, SelectInputFile, ValidateKey} from '../wailsjs/go/main/App';
import './App.css';

function App() {
    const [operation, setOperation] = useState('encrypt');
    const [inputFile, setInputFile] = useState('');
    const [key, setKey] = useState('');
    const [result, setResult] = useState(null);
    const [keyInfo, setKeyInfo] = useState(null);
    const [loading, setLoading] = useState(false);

    // Обработка выбора файла
    const handleSelectFile = async () => {
        try {
            const file = await SelectInputFile();
            if (file) {
                setInputFile(file);
                setResult(null);
            }
        } catch (err) {
            console.error('Ошибка выбора файла:', err);
        }
    };

    // Валидация ключа при вводе
    const handleKeyChange = async (value) => {
        setKey(value);
        if (value) {
            try {
                const info = await ValidateKey(value);
                setKeyInfo(info);
            } catch (err) {
                setKeyInfo(null);
            }
        } else {
            setKeyInfo(null);
        }
    };

    // Обработка операции
    const handleProcess = async () => {
        if (!key) {
            setResult({success: false, message: 'Введите ключ'});
            return;
        }

        if (!inputFile) {
            setResult({success: false, message: 'Выберите файл'});
            return;
        }

        // Предварительная проверка ключа
        const keyValidation = await ValidateKey(key);
        if (!keyValidation.valid) {
            setResult({success: false, message: keyValidation.message});
            return;
        }

        setLoading(true);
        setResult(null);

        try {
            let res;
            if (operation === 'encrypt') {
                res = await EncryptFile(inputFile, key);
            } else {
                res = await DecryptFile(inputFile, key);
            }
            setResult(res);
        } catch (err) {
            setResult({success: false, message: `Ошибка: ${err}`});
        }

        setLoading(false);
    };

    return (
        <div className="app-container">
            <header className="app-header">
                <h1>🔐 LFSR Шифрование</h1>
                <p className="subtitle">Полином: x³⁷ + x¹² + x¹⁰ + x² + 1</p>
            </header>

            <main className="main-content">
                {/* Переключатель операции */}
                <div className="toggle-section">
                    <div className="toggle-group">
                        <label>Операция:</label>
                        <div className="toggle-buttons">
                            <button
                                className={`toggle-btn ${operation === 'encrypt' ? 'active' : ''}`}
                                onClick={() => {
                                    setOperation('encrypt');
                                    setResult(null);
                                }}
                            >
                                🔒 Шифрование
                            </button>
                            <button
                                className={`toggle-btn ${operation === 'decrypt' ? 'active' : ''}`}
                                onClick={() => {
                                    setOperation('decrypt');
                                    setResult(null);
                                }}
                            >
                                🔓 Дешифрование
                            </button>
                        </div>
                    </div>
                </div>

                {/* Ввод данных */}
                <div className="input-section">
                    {/* Выбор файла */}
                    <div className="file-input-group">
                        <label>Входной файл:</label>
                        <div className="file-input-row">
                            <input
                                type="text"
                                value={inputFile}
                                readOnly
                                placeholder="Файл не выбран..."
                                className="file-path-input"
                            />
                            <button onClick={handleSelectFile} className="select-file-btn">
                                📁 Выбрать файл
                            </button>
                        </div>
                        <p className="input-hint">
                            Поддерживаются все типы файлов: изображения, видео, аудио, документы, архивы и др.
                        </p>
                    </div>

                    {/* Ввод ключа */}
                    <div className="key-input-group">
                        <label>Ключ шифрования (37 бит):</label>
                        <input
                            type="text"
                            value={key}
                            onChange={(e) => handleKeyChange(e.target.value)}
                            placeholder="Введите 37 символов (только 0 и 1)..."
                            className={`key-input ${keyInfo ? (keyInfo.valid ? 'valid' : 'invalid') : ''}`}
                            maxLength={50}
                        />
                        <div className="key-counter">
              <span className={key.length === 37 ? 'valid-count' : key.length > 37 ? 'invalid-count' : ''}>
                {key.length}
              </span> / 37 символов
                        </div>
                    </div>

                    {/* Информация о ключе */}
                    {keyInfo && (
                        <div className={`key-info ${keyInfo.valid ? 'valid' : 'invalid'}`}>
                            <div className="key-info-header">
                                {keyInfo.valid ? '✅' : '❌'} {keyInfo.message}
                            </div>
                            {key && (
                                <div className="binary-key-display">
                                    <strong>Введённый ключ:</strong>
                                    <div className="binary-value">
                                        {key.split('').map((char, index) => (
                                            <span
                                                key={index}
                                                className={char === '0' || char === '1' ? 'valid-char' : 'invalid-char'}
                                            >
                        {char}
                      </span>
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>
                    )}

                    {/* Кнопка обработки */}
                    <button
                        onClick={handleProcess}
                        className="process-btn"
                        disabled={loading || !keyInfo?.valid || !inputFile}
                    >
                        {loading
                            ? '⏳ Обработка...'
                            : (operation === 'encrypt' ? '🔒 Зашифровать файл' : '🔓 Расшифровать файл')
                        }
                    </button>
                </div>

                {/* Результат */}
                {result && (
                    <div className={`result-section ${result.success ? 'success' : 'error'}`}>
                        <h3>{result.success ? '✅ Успешно' : '❌ Ошибка'}</h3>

                        {result.success ? (
                            <>
                                <div className="result-info">
                                    <p><strong>Статус:</strong> {result.message}</p>
                                    <p><strong>Сохранено в:</strong> <span
                                        className="file-path">{result.outputPath}</span></p>
                                </div>

                                <div className="crypto-info">
                                    <h4>🔑 Криптографическая информация</h4>

                                    <div className="info-block">
                                        <strong>Использованный ключ (37 бит):</strong>
                                        <div className="binary-display">{result.binaryKey}</div>
                                    </div>

                                    <div className="info-block">
                                        <strong>Сгенерированный keystream (первые 256 бит):</strong>
                                        <div className="binary-display">{result.keyStream}</div>
                                    </div>

                                </div>
                            </>
                        ) : (
                            <p className="error-message">{result.message}</p>
                        )}
                    </div>
                )}
            </main>

        </div>
    );
}

export default App;