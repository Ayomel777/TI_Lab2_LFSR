import { useState, useEffect } from 'react';
import {
    SelectAndReadFile,
    Encrypt,
    Decrypt,
    ValidateKey,
    ValidateInput,
    SaveToFile
} from '../wailsjs/go/main/App';
import './App.css';

function App() {
    const [mode, setMode] = useState('encrypt');
    const [inputText, setInputText] = useState('');
    const [register, setRegister] = useState('');
    const [keyStream, setKeyStream] = useState('');
    const [outputText, setOutputText] = useState('');
    const [registerInfo, setRegisterInfo] = useState(null);
    const [inputInfo, setInputInfo] = useState(null);
    const [status, setStatus] = useState({ type: '', message: '' });
    const [loading, setLoading] = useState(false);
    const [canSave, setCanSave] = useState(false);
    const [fileExtension, setFileExtension] = useState('bin');

    useEffect(() => {
        const validate = async () => {
            if (register) {
                const info = await ValidateKey(register);
                setRegisterInfo(info);
            } else {
                setRegisterInfo(null);
            }
        };
        validate();
    }, [register]);

    useEffect(() => {
        const validate = async () => {
            if (inputText) {
                const info = await ValidateInput(inputText);
                setInputInfo(info);
            } else {
                setInputInfo(null);
            }
        };
        validate();
    }, [inputText]);

    const clearOutput = () => {
        setOutputText('');
        setKeyStream('');
        setCanSave(false);
        setStatus({ type: '', message: '' });
    };

    const getInputLabel = () => mode === 'encrypt' ? 'Исходный текст (биты):' : 'Шифротекст (биты):';
    const getOutputLabel = () => mode === 'encrypt' ? 'Шифротекст (биты):' : 'Расшифрованный текст (биты):';
    const getActionLabel = () => mode === 'encrypt' ? 'Зашифровать' : 'Расшифровать';

    const handleOpenFile = async () => {
        try {
            const fileResult = await SelectAndReadFile();
            if (fileResult.success) {
                setInputText(fileResult.bits);
                const fileName = fileResult.filePath.split(/[/\\]/).pop();
                const ext = fileName.includes('.') ? fileName.split('.').pop() : 'bin';
                setFileExtension(ext);
                setStatus({ type: 'success', message: `Файл загружен: ${fileName} (${fileResult.fileSize} байт = ${fileResult.bits.length} бит)` });
                clearOutput();
            } else if (fileResult.message !== "Файл не выбран") {
                setStatus({ type: 'error', message: fileResult.message });
            }
        } catch (err) {
            setStatus({ type: 'error', message: `Ошибка: ${err}` });
        }
    };

    const handleSaveFile = async () => {
        if (!outputText) {
            setStatus({ type: 'error', message: 'Нет данных для сохранения' });
            return;
        }

        try {
            const prefix = mode === 'encrypt' ? 'encrypted' : 'decrypted';
            const defaultName = `${prefix}.${fileExtension}`;
            const res = await SaveToFile(outputText, defaultName);
            if (res.success) {
                setStatus({ type: 'success', message: res.message });
            } else {
                setStatus({ type: 'error', message: res.message });
            }
        } catch (err) {
            setStatus({ type: 'error', message: `Ошибка: ${err}` });
        }
    };

    const handleAction = async () => {
        if (!inputInfo?.valid) {
            setStatus({ type: 'error', message: inputInfo?.message || 'Введите данные' });
            return;
        }
        if (!registerInfo?.valid) {
            setStatus({ type: 'error', message: registerInfo?.message || 'Введите корректный регистр' });
            return;
        }

        setLoading(true);
        setStatus({ type: '', message: '' });

        try {
            let res;
            if (mode === 'encrypt') {
                res = await Encrypt(inputText, register);
            } else {
                res = await Decrypt(inputText, register);
            }

            if (res.success) {
                setOutputText(res.cipherText);
                setKeyStream(res.keyStream);
                setCanSave(true);
                setStatus({
                    type: 'success',
                    message: `${mode === 'encrypt' ? 'Зашифровано' : 'Расшифровано'} ${res.bitsCount} бит`
                });
            } else {
                setStatus({ type: 'error', message: res.message });
            }
        } catch (err) {
            setStatus({ type: 'error', message: `Ошибка: ${err}` });
        }

        setLoading(false);
    };

    const handleModeChange = (newMode) => {
        setMode(newMode);
        clearOutput();
    };

    const handleRegisterChange = (value) => {
        setRegister(value);
        clearOutput();
    };

    const handleInputChange = (value) => {
        setInputText(value);
        setFileExtension('bin');
        clearOutput();
    };

    const canExecute = inputInfo?.valid && registerInfo?.valid && !loading;

    return (
        <div className="app-container">
            <header className="app-header">
                <h1>Лабораторная работа — LFSR Шифрование</h1>
                <p className="subtitle">Полином: x³⁷ + x¹² + x¹⁰ + x² + 1</p>
            </header>

            <div className="app-content">
                <div className="mode-selector">
                    <label className={`mode-option ${mode === 'encrypt' ? 'active' : ''}`}>
                        <input
                            type="radio"
                            name="mode"
                            checked={mode === 'encrypt'}
                            onChange={() => handleModeChange('encrypt')}
                        />
                        <span>Шифрование</span>
                    </label>
                    <label className={`mode-option ${mode === 'decrypt' ? 'active' : ''}`}>
                        <input
                            type="radio"
                            name="mode"
                            checked={mode === 'decrypt'}
                            onChange={() => handleModeChange('decrypt')}
                        />
                        <span>Дешифрование</span>
                    </label>
                </div>

                <div className="file-buttons">
                    <button onClick={handleOpenFile} className="file-btn open-btn">
                        📁 Открыть файл
                    </button>
                    <button
                        onClick={handleSaveFile}
                        className="file-btn save-btn"
                        disabled={!canSave}
                    >
                        💾 Сохранить результат
                    </button>
                </div>

                <div className="register-section">
                    <div className="register-header">
                        <span className="register-label">Начальное состояние регистра (ровно 37 бит):</span>
                        <span className={`register-count ${register.length === 37 ? 'valid' : register.length > 37 ? 'invalid' : ''}`}>
              {register.length}
            </span>
                        <span className="register-total">/ 37</span>
                    </div>
                    <input
                        type="text"
                        value={register}
                        onChange={(e) => handleRegisterChange(e.target.value)}
                        placeholder="Введите 37 бит (только 0 и 1)..."
                        className={`register-input ${registerInfo ? (registerInfo.valid ? 'valid' : 'invalid') : ''}`}
                        maxLength={37}
                    />
                    {registerInfo && !registerInfo.valid && (
                        <div className="register-error">{registerInfo.message}</div>
                    )}
                </div>

                <div className="columns-container">
                    <div className="column">
                        <div className="column-header">
                            <span className="column-label">{getInputLabel()}</span>
                            {inputInfo && (
                                <span className={`column-count ${inputInfo.valid ? 'valid' : 'invalid'}`}>
                  {inputInfo.valid ? `${inputInfo.bitsCount} бит` : '⚠️'}
                </span>
                            )}
                        </div>
                        <textarea
                            value={inputText}
                            onChange={(e) => handleInputChange(e.target.value)}
                            placeholder="Введите биты (только 0 и 1) или загрузите файл..."
                            className="column-textarea"
                        />
                        {inputInfo && !inputInfo.valid && (
                            <div className="column-error">{inputInfo.message}</div>
                        )}
                    </div>

                    <div className="column">
                        <div className="column-header">
                            <span className="column-label">Биты ключа (keystream):</span>
                            {keyStream && (
                                <span className="column-count valid">{keyStream.replace('...', '').length}+ бит</span>
                            )}
                        </div>
                        <textarea
                            value={keyStream}
                            readOnly
                            placeholder="Здесь появятся биты ключа после обработки..."
                            className="column-textarea readonly"
                        />
                    </div>

                    <div className="column">
                        <div className="column-header">
                            <span className="column-label">{getOutputLabel()}</span>
                            {outputText && (
                                <span className="column-count valid">{outputText.length} бит</span>
                            )}
                        </div>
                        <textarea
                            value={outputText}
                            readOnly
                            placeholder="Здесь появится результат после обработки..."
                            className="column-textarea readonly"
                        />
                    </div>
                </div>

                <button
                    onClick={handleAction}
                    className={`action-btn ${mode === 'encrypt' ? 'encrypt' : 'decrypt'}`}
                    disabled={!canExecute}
                >
                    {loading ? '⏳ Обработка...' : `${mode === 'encrypt' ? '🔒' : '🔓'} ${getActionLabel()}`}
                </button>

                {status.message && (
                    <div className={`status-bar ${status.type}`}>
                        {status.type === 'success' ? '✅' : '❌'} {status.message}
                    </div>
                )}
            </div>

        </div>
    );
}

export default App;