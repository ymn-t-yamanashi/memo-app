import React, { useState, useEffect } from 'react';
import axios from 'axios';

const API_URL = 'http://localhost:8080/memos';

function App() {
  const [memos, setMemos] = useState([]);
  const [text, setText] = useState('');

  // 1. データ取得 (Read)
  const fetchMemos = async () => {
    const res = await axios.get(API_URL);
    setMemos(res.data || []);
  };

  useEffect(() => { fetchMemos(); }, []);

  // 2. データ保存 (Create)
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!text.trim()) return;
    await axios.post(API_URL, { body: text });
    setText('');
    fetchMemos();
  };

  return (
    <div style={{ maxWidth: '500px', margin: '40px auto', fontFamily: 'sans-serif' }}>
      <h1>Memo App</h1>
      <form onSubmit={handleSubmit} style={{ display: 'flex', gap: '10px' }}>
        <input 
          style={{ flex: 1, padding: '10px' }}
          value={text} 
          onChange={(e) => setText(e.target.value)} 
          placeholder="今日の一言を記録..."
        />
        <button style={{ padding: '10px 20px' }} type="submit">保存</button>
      </form>

      <div style={{ marginTop: '20px' }}>
        {memos.map(m => (
          <div key={m.id} style={{ padding: '10px', borderBottom: '1px solid #eee' }}>
            {m.body}
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;