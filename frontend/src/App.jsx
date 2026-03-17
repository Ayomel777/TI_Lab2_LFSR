import { useState } from 'react';
import './App.css';

function App() {
    const [key, setKey] = useState('');
    const [inputFile, setInputFile] = useState('');
    const [outputFile, setOutputFile] = useState('');
    const [status, setStatus] = useState('');
    const [loading, setLoading] = useState(false);

    // Установка ключа
    const handleSetKey = async () => {
        if (!/^[01]{37}$/.test(key)) {
            setStatus('❌ Ключ должен содержать ровно 37 символов (только 0 и 1)');
            return;
        }
        setLoading(true);
        try {
            await window.go.main.App.SetKey(key);
            const info = await window.go.main.App.GetKeyInfo();
            setKey(info);
            setStatus('✅ Ключ успешно установлен');
        } catch (err) {
            setStatus(`❌ Ошибка: ${err}`);
        } finally {
            setLoading(false);
        }
    };

    // Выбор входного файла
    const handleSelectInput = async () => {
        try {
            const result = await window.go.main.App.OpenFileDialog();
            if (result) setInputFile(result);
        } catch (err) {
            setStatus(`❌ Ошибка при выборе файла: ${err}`);
        }
    };

    // Выбор выходного файла
    const handleSelectOutput = async () => {
        try {
            const result = await window.go.main.App.SaveFileDialog();
            if (result) setOutputFile(result);
        } catch (err) {
            setStatus(`❌ Ошибка при выборе файла: ${err}`);
        }
    };

    // Шифрование
    const handleEncrypt = async () => {
        if (!inputFile || !outputFile) {
            setStatus('❌ Выберите входной и выходной файлы');
            return;
        }
        if (!key) {
            setStatus('❌ Сначала установите ключ');
            return;
        }
        setLoading(true);
        try {
            await window.go.main.App.Encrypt(inputFile, outputFile);
            setStatus('✅ Шифрование завершено');
        } catch (err) {
            setStatus(`❌ Ошибка шифрования: ${err}`);
        } finally {
            setLoading(false);
        }
    };

    // Дешифрование
    const handleDecrypt = async () => {
        if (!inputFile || !outputFile) {
            setStatus('❌ Выберите входной и выходной файлы');
            return;
        }
        if (!key) {
            setStatus('❌ Сначала установите ключ');
            return;
        }
        setLoading(true);
        try {
            await window.go.main.App.Decrypt(inputFile, outputFile);
            setStatus('✅ Дешифрование завершено');
        } catch (err) {
            setStatus(`❌ Ошибка дешифрования: ${err}`);
        } finally {
            setLoading(false);
        }
    };

    // Очистка всех полей
    const handleClear = () => {
        setKey('');
        setInputFile('');
        setOutputFile('');
        setStatus('');
    };

    return (
        <div className="app">
            <h1>🔐 LFSR File Encryptor (степень 37)</h1>
            <p className="subtitle">
                Многочлен: x<sup>37</sup> + x<sup>6</sup> + x<sup>4</sup> + x + 1
            </p>

            {/* Ключ */}
            <div className="row">
                <input
                    type="text"
                    value={key}
                    onChange={(e) => setKey(e.target.value)}
                    placeholder="37-битный двоичный ключ"
                    maxLength="37"
                    pattern="[01]{37}"
                    disabled={loading}
                    className="keyword"
                />
                <button onClick={handleSetKey} disabled={loading} className="btn setkey">
                    Установить ключ
                </button>
            </div>

            {/* Входной файл */}
            <div className="row file-row">
                <input
                    type="text"
                    value={inputFile}
                    readOnly
                    placeholder="Входной файл не выбран"
                    className="file-path"
                />
                <button onClick={handleSelectInput} disabled={loading} className="btn file">
                    📂 Обзор
                </button>
            </div>

            {/* Выходной файл */}
            <div className="row file-row">
                <input
                    type="text"
                    value={outputFile}
                    readOnly
                    placeholder="Выходной файл не выбран"
                    className="file-path"
                />
                <button onClick={handleSelectOutput} disabled={loading} className="btn save">
                    💾 Сохранить как
                </button>
            </div>

            {/* Кнопки действий */}
            <div className="buttons">
                <button onClick={handleEncrypt} disabled={loading} className="btn encrypt">
                    ЗАШИФРОВАТЬ
                </button>
                <button onClick={handleDecrypt} disabled={loading} className="btn decrypt">
                    РАСШИФРОВАТЬ
                </button>
                <button onClick={handleClear} disabled={loading} className="btn clear">
                    🗑️ ОЧИСТИТЬ
                </button>
            </div>

            {/* Статус */}
            {status && (
                <div className="status">
                    {status}
                </div>
            )}

            {/* Индикатор загрузки */}
            {loading && (
                <div className="loading-overlay">
                    <div className="spinner"></div>
                    <p>Обработка...</p>
                </div>
            )}
        </div>
    );
}

export default App;