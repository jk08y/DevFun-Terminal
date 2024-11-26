import React, { useState, useRef, useEffect } from 'react';

const EnhancedTerminal = () => {
  const [commands, setCommands] = useState([]);
  const [input, setInput] = useState('');
  const [commandHistory, setCommandHistory] = useState([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const [isFullScreen, setIsFullScreen] = useState(false);
  const bottomRef = useRef(null);
  const inputRef = useRef(null);

  // Comprehensive command responses with more context and humor
  const funnyResponses = {
    'relationship status': {
      response: 'Committed to main branch 💍\n- Relationship type: Git merge\n- Conflict resolution: Always communicate',
      type: 'success'
    },
    'coffee level': {
      response: 'Exception: Coffee Not Found ☕️\n- Current caffeine status: Critical\n- Recommended action: Brew immediately',
      type: 'warning'
    },
    'system info': {
      response: 'System Specs:\n- OS: Human v3.0\n- Processor: Brain 9000\n- RAM: Coffee-powered\n- Storage: Memories & Memes',
      type: 'info'
    },
    'motivate': {
      response: 'Motivation Booster Activated 🚀\n- You are a coding wizard\n- Every bug is just an undocumented feature\n- Keep pushing your limits (and your code)',
      type: 'success'
    },
    'joke': {
      response: 'Why do programmers prefer dark mode?\nBecause light attracts bugs! 🐛😂',
      type: 'fun'
    },
    'date': {
      response: `Current Date & Time: ${new Date().toLocaleString()}`,
      type: 'info'
    },
    'weather': {
      response: 'Forecast: 100% chance of code with scattered bugs 🌦️\n- Temperature: Hot debug session\n- Humidity: Sweaty keyboard',
      type: 'info'
    },
    // New commands from the provided list
    'debug life': {
      response: 'Found 99 problems, but a bug ain\'t one 🐛\n- Life debugging status: In progress\n- Error handling: Minimal',
      type: 'success'
    },
    'work status': {
      response: 'Current mission: git push --force-with-coffee ⚡️\n- Productivity: Caffeine-driven\n- Commits: Highly caffeinated',
      type: 'info'
    },
    'weekend plans': {
      response: 'Current algorithm:\nwhile(true) { sleep(); }\n- Rest mode: Activated\n- Productivity: Paused',
      type: 'fun'
    },
    'current task': {
      response: 'Task Status: Turning coffee into code... Loading... ⌛️\n- Progress: Caffeine to Code Conversion\n- ETA: Next coffee break',
      type: 'warning'
    },
    'ping heart': {
      response: 'Connection Status: ❤️\n- Signal Strength: Crushingly Strong\n- Latency: Instant Butterflies',
      type: 'success'
    },
    'git status': {
      response: 'Branch Analysis:\n- Current Branch: Life\n- Ahead of master by: 42 commits 🌟\n- Merge conflicts: Minimal',
      type: 'info'
    },
    'sudo make coffee': {
      response: 'Permission Denied 🚫\n- Only authorized baristas can perform this action\n- Alternative: Manual coffee brewing recommended ☕️',
      type: 'error'
    },
    'ls friends': {
      response: 'Friend Directory:\n- Contents: []\n- Error: No friends found 😢\n- Suggestion: Debug social skills',
      type: 'warning'
    },
    'rm -rf stress': {
      response: 'Error: stress is a read-only file in your life system 🛡️\n- Recommended Action: sudo vacation\n- Stress Protection: Enabled',
      type: 'error'
    },
    'cat wisdom': {
      response: 'Wisdom Kernel Activated 🐱‍💻\n- Rule #1: If code compiles, don\'t touch it\n- Life Optimization: Minimal Interference',
      type: 'success'
    },
    'whoami': {
      response: 'User Profile:\n- Role: Coding Ninja 💻\n- Special Abilities:\n  * Caffeine Resistance\n  * Bug Annihilation\n  * Stack Overflow Fluency',
      type: 'info'
    }
  }



  funnyResponses['help'] = {
    response: 'Available Commands:\n' + Object.keys(funnyResponses)
      .filter(cmd => cmd !== 'help' && cmd !== 'clear')
      .map(cmd => `• ${cmd}`)
      .join('\n'),
    type: 'info'
  };

  // Enhanced command handling with more robust processing
  const handleCommand = (e) => {
    e.preventDefault();
    const trimmedInput = input.trim().toLowerCase();
    
    if (trimmedInput === '') return;

    // Add to command history
    const updatedHistory = [...commandHistory, trimmedInput];
    setCommandHistory(updatedHistory);
    setHistoryIndex(updatedHistory.length);

    // Find response
    const commandResponse = funnyResponses[trimmedInput] || {
      response: `Command not found: "${input}". Try "help" for available commands 🤔`,
      type: 'error'
    };

    // Special commands
    if (trimmedInput === 'clear') {
      setCommands([]);
      return;
    }

    // Add command to terminal output
    handleAddCommand({
      type: 'input',
      text: input
    });
    handleAddCommand({
      type: 'output',
      text: commandResponse.response,
      style: commandResponse.type || 'default'
    });

    // Reset input
    setInput('');
  };

  // Improved command addition with styling
  const handleAddCommand = (command) => {
    setCommands(prevCommands => [...prevCommands, command]);
  };

  // Keyboard navigation for command history
  const handleKeyDown = (e) => {
    if (e.key === 'ArrowUp') {
      e.preventDefault();
      if (historyIndex > 0) {
        const newIndex = historyIndex - 1;
        setInput(commandHistory[newIndex]);
        setHistoryIndex(newIndex);
      }
    } else if (e.key === 'ArrowDown') {
      e.preventDefault();
      if (historyIndex < commandHistory.length - 1) {
        const newIndex = historyIndex + 1;
        setInput(commandHistory[newIndex]);
        setHistoryIndex(newIndex);
      } else {
        setInput('');
        setHistoryIndex(commandHistory.length);
      }
    }
  };

  // Auto-scroll to bottom
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [commands]);

  // Focus input on mount and after commands
  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  return (
    <div className={`min-h-screen bg-gray-900 flex items-center justify-center p-4 transition-all duration-300 ${isFullScreen ? 'fixed inset-0 z-50' : ''}`}>
      <div className={`w-full ${isFullScreen ? 'max-w-full h-full' : 'max-w-2xl'} bg-gray-800 rounded-lg shadow-xl overflow-hidden`}>
        {/* Terminal Header */}
        <div className="p-2 bg-gray-900 flex items-center justify-between">
          <div className="flex gap-2">
            <div className="w-3 h-3 rounded-full bg-red-500"></div>
            <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
            <div className="w-3 h-3 rounded-full bg-green-500"></div>
          </div>
          <div className="text-center flex-1">
            <span className="text-gray-400 font-mono">jk08y@pro-terminal</span>
          </div>
          <div className="flex items-center gap-2">
            <button 
              onClick={() => setIsFullScreen(!isFullScreen)} 
              className="text-gray-400 hover:text-white"
            >
              {isFullScreen ? '🗝️' : '🖥️'}
            </button>
          </div>
        </div>
        
        {/* Terminal Content */}
        <div className={`p-4 overflow-y-auto font-mono text-sm ${isFullScreen ? 'h-[calc(100vh-100px)]' : 'h-96'}`}>
          <div className="text-green-400 mb-4">
            Welcome to DevFun Terminal v2.0.0 🚀 Type 'help' for available commands
          </div>
          
          {/* Command History Rendering */}
          {commands.map((cmd, index) => (
            <div 
              key={index} 
              className={`mb-2 ${
                cmd.style === 'error' ? 'text-red-400' : 
                cmd.style === 'warning' ? 'text-yellow-400' : 
                cmd.style === 'success' ? 'text-green-400' : 
                cmd.style === 'info' ? 'text-blue-400' : 
                cmd.style === 'fun' ? 'text-purple-400' : 
                'text-white'
              }`}
            >
              {cmd.type === 'input' ? (
                <div>
                  <span className="text-pink-400">➜</span> {cmd.text}
                </div>
              ) : (
                <div className="pl-4 whitespace-pre-wrap">{cmd.text}</div>
              )}
            </div>
          ))}
          
          {/* Command Input */}
          <form onSubmit={handleCommand} className="flex items-center">
            <span className="text-pink-400">➜</span>
            <input
              ref={inputRef}
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              className="flex-1 ml-2 bg-transparent text-white outline-none"
              autoFocus
              placeholder="Type a command (try 'help')"
            />
          </form>
          <div ref={bottomRef} />
        </div>
      </div>
    </div>
  );
};

export default EnhancedTerminal;
