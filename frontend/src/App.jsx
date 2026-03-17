import { useState } from 'react';
import {
    SelectInputFile,
    EncryptFile,
    DecryptFile,
    EncryptText,
    DecryptText,
    ConvertToBinaryKey
} from '../wailsjs/go/main/App';
import './App.css';

function App() {
    const [mode, setMode] = useState('file'); // 'file' или 'text'
    const [operation, setOperation] = useState('encrypt'); // 'encrypt' или 'decrypt'
    const [inputFile, setInputFile] = useState('');
    const [inputText, setInputText] = useState('');
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

    // Конвертация и отображение ключа
    const handleKeyChange = async (value) => {
        setKey(value);
        if (value) {
            try {
                const info = await ConvertToBinaryKey(value);
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
            setResult({ success: false, message: 'Введите ключ' });
            return;
        }

        setLoading(true);
        setResult(null);

        try {
            let res;
            if (mode === 'file') {
                if (!inputFile) {
                    setResult({ success: false, message: 'Выберите файл' });
                    setLoading(false);
                    return;
                }
                if (operation === 'encrypt') {
                    res = await EncryptFile(inputFile, key);
                } else {
                    res = await DecryptFile(inputFile, key);
                }
            } else {
                if (!inputText) {
                    setResult({ success: false, message: 'Введите текст' });
                    setLoading(false);
                    return;
                }
                if (operation === 'encrypt') {
                    res = await EncryptText(inputText, key);
                } else {
                    res = await DecryptText(inputText, key);
                }
            }
            setResult(res);
        } catch (err) {
            setResult({ success: false, message: `Ошибка: ${err}` });
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
                {/* Переключатели режима */}
                <div className="toggle-section">
                    <div className="toggle-group">
                        <label>Режим:</label>
                        <div className="toggle-buttons">
                            <button
                                className={`toggle-btn ${mode === 'file' ? 'active' : ''}`}
                                onClick={() => { setMode('file'); setResult(null); }}
                            >
                                📁 Файл
                            </button>
                            <button
                                className={`toggle-btn ${mode === 'text' ? 'active' : ''}`}
                                onClick={() => { setMode('text'); setResult(null); }}
                            >
                                📝 Текст
                            </button>
                        </div>
                    </div>

                    <div className="toggle-group">
                        <label>Операция:</label>
                        <div className="toggle-buttons">
                            <button
                                className={`toggle-btn ${operation === 'encrypt' ? 'active' : ''}`}
                                onClick={() => { setOperation('encrypt'); setResult(null); }}
                            >
                                🔒 Шифрование
                            </button>
                            <button
                                className={`toggle-btn ${operation === 'decrypt' ? 'active' : ''}`}
                                onClick={() => { setOperation('decrypt'); setResult(null); }}
                            >
                                🔓 Дешифрование
                            </button>
                        </div>
                    </div>
                </div>

                {/* Ввод данных */}
                <div className="input-section">
                    {mode === 'file' ? (
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
                                    Выбрать файл
                                </button>
                            </div>
                        </div>
                    ) : (
                        <div className="text-input-group">
                            <label>
                                {operation === 'encrypt' ? 'Текст для шифрования:' : 'Зашифрованный текст (Base64):'}
                            </label>
                            <textarea
                                value={inputText}
                                onChange={(e) => setInputText(e.target.value)}
                                placeholder={operation === 'encrypt'
                                    ? 'Введите текст для шифрования...'
                                    : 'Вставьте зашифрованный текст в формате Base64...'}
                                className="text-input"
                                rows={4}
                            />
                        </div>
                    )}

                    {/* Ввод ключа */}
                    <div className="key-input-group">
                        <label>Ключ шифрования:</label>
                        <input
                            type="text"
                            value={key}
                            onChange={(e) => handleKeyChange(e.target.value)}
                            placeholder="Введите ключ (любые символы или бинарный)..."
                            className="key-input"
                        />
                        <p className="key-hint">
                            Ключ будет автоматически преобразован в бинарный формат.
                            Минимальная длина: 37 бит.
                        </p>
                    </div>

                    {/* Информация о ключе */}
                    {keyInfo && (
                        <div className={`key-info ${keyInfo.valid ? 'valid' : 'invalid'}`}>
                            <h4>Информация о ключе:</h4>
                            <p><strong>Статус:</strong> {keyInfo.message}</p>
                            {keyInfo.valid && (
                                <>
                                    <p><strong>Длина:</strong> {keyInfo.keyLength} бит</p>
                                    <div className="binary-key-display">
                                        <strong>Бинарное представление:</strong>
                                        <div className="binary-value">
                                            {keyInfo.binaryKey.length > 100
                                                ? keyInfo.binaryKey.substring(0, 100) + '...'
                                                : keyInfo.binaryKey}
                                        </div>
                                    </div>
                                </>
                            )}
                        </div>
                    )}

                    {/* Кнопка обработки */}
                    <button
                        onClick={handleProcess}
                        className="process-btn"
                        disabled={loading}
                    >
                        {loading ? '⏳ Обработка...' : (operation === 'encrypt' ? '🔒 Зашифровать' : '🔓 Расшифровать')}
                    </button>
                </div>

                {/* Результат */}
                {result && (
                    <div className={`result-section ${result.success ? 'success' : 'error'}`}>
                        <h3>{result.success ? '✅ Успех' : '❌ Ошибка'}</h3>

                        {result.success ? (
                            <>
                                {mode === 'file' ? (
                                    <div className="result-info">
                                        <p><strong>Сообщение:</strong> {result.message}</p>
                                        <p><strong>Путь к файлу:</strong> {result.outputPath}</p>
                                        <p><strong>Размер исходных данных:</strong> {result.originalSize} байт</p>
                                        <p><strong>Размер обработанных данных:</strong> {result.processedSize} байт</p>
                                    </div>
                                ) : (
                                    <div className="result-info">
                                        <p><strong>Результат:</strong></p>
                                        <div className="result-text">
                                            {result.message}
                                        </div>
                                        <p><strong>Размер:</strong> {result.processedSize} байт</p>
                                    </div>
                                )}

                                <div className="crypto-info">
                                    <h4>🔑 Криптографическая информация:</h4>

                                    <div className="info-block">
                                        <strong>Бинарный ключ (первые 512 бит):</strong>
                                        <div className="binary-display">{result.binaryKey}</div>
                                    </div>

                                    <div className="info-block">
                                        <strong>Сгенерированный keystream (первые 256 бит):</strong>
                                        <div className="binary-display">{result.keyStream}</div>
                                    </div>

                                    <div className="symmetry-note">
                                        <p>ℹ️ <strong>Симметричное шифрование:</strong> Операции шифрования и дешифрования
                                            идентичны (XOR с одинаковым keystream). Для расшифровки используйте тот же ключ.</p>
                                    </div>
                                </div>
                            </>
                        ) : (
                            <p className="error-message">{result.message}</p>
                        )}
                    </div>
                )}
            </main>

            <footer className="app-footer">
                <p>LFSR (Linear Feedback Shift Register) с полиномом степени 37</p>
            </footer>
        </div>
    );
}

export default App;