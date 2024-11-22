import React, { useState, useRef, useEffect } from 'react';

const Terminal = () => {
  const [commands, setCommands] = useState([]);
  const [input, setInput] = useState('');
  const bottomRef = useRef(null);

  const funnyResponses = {
    'relationship status': 'Committed to main branch 💍',
    'coffee level': 'Exception: Coffee Not Found ☕️',
    'debug life': 'Found 99 problems, but a bug ain\'t one 🐛',
    'work status': 'git push --force-with-coffee ⚡️',
    'motivation': 'Running on caffeine and dreams... and Stack Overflow 🚀',
    'weekend plans': 'while(true) { sleep(); }',
    'current task': 'Turning coffee into code... Loading... ⌛️',
    'ping heart': 'Connection established with crush <3',
    'git status': 'Your branch is ahead of life by 42 commits 🌟',
    'sudo make coffee': 'Permission denied: Only authorized baristas can perform this action ☕️',
    'ls friends': 'No friends found 😢',
    'rm -rf stress': 'Error: stress is a read-only file in your life system. Try "sudo vacation" instead 🌴',
    'cat wisdom': 'If code compiles, don\'t touch it. 🐱‍💻',
    'apt-get happiness': 'Package not found. Please try "sudo self-care" 🌱',
    'man woman': 'Error: man pages do not contain enough information on woman. Try asking nicely 💃',
    'ping boss': 'Request timed out. Boss is unreachable 😎',
    'top': 'Currently using 99% of brain capacity on coffee processing ☕️',
    'whoami': 'You are a coding ninja 💻',
    'cp motivation /dev/me': 'Success: Motivation copied! Boosting productivity by 20% 🚀',
    'df -h': 'Emotional space: 10% free. Recommend backing up memories 📦',
    'ls': 'Files: ./overwhelming_deadlines ./mild_panic ./restarting_project_(again)',
    'grep "love" life': 'Love not found. Check your relationships file 🥀',
    'uname -a': 'System: Human 2.0, Powered by caffeine, Random reboots due to lack of sleep 🧠',
    'sudo find happiness': 'Permission denied: Happiness must be self-discovered 🧘‍♂️',
    'chmod +x productivity': 'Productivity level set to executable! 🎉',
    'cd /success': 'No such file or directory found. Try "/hardwork" instead 💪',
    'history -c': 'Nice try, but memories cannot be deleted that easily 😆',
    'ps -ef | grep "peace"': 'Process "peace" not found. Try meditation 🧘',
    'chmod 777 everything': 'Warning: Security team has been notified. Your future self disapproves 👀',
    'kill -9 laziness': 'Error: laziness is an unkillable process. Please try again after coffee ☕️',
    'echo "sleep"': 'sleep: command not found. Try "insomnia --force" instead 💤',
    'help': 'Available commands: relationship status, coffee level, debug life, work status, motivation, weekend plans, current task, ping heart, git status, sudo make coffee, ls friends, rm -rf stress, cat wisdom, apt-get happiness, man woman, ping boss, top, whoami, cp motivation /dev/me, df -h, grep "love" life, uname -a, sudo find happiness, chmod +x productivity, cd /success, history -c, ps -ef | grep "peace", kill -9 laziness, echo "sleep"',
    'clear': 'CLEAR_COMMAND'
  };

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [commands]);

  const handleCommand = (e) => {
    e.preventDefault();
    if (input.trim() === '') return;

    const response = funnyResponses[input.toLowerCase()] || 'Command not found: Try "help" for available commands 🤔';
    
    if (input.toLowerCase() === 'clear') {
      setCommands([]);
    } else {
      setCommands([...commands, { type: 'input', text: input }, { type: 'output', text: response }]);
    }
    setInput('');
  };

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <div className="w-full max-w-2xl bg-gray-800 rounded-lg shadow-xl overflow-hidden">
        <div className="p-2 bg-gray-900 flex items-center gap-2">
          <div className="flex gap-2">
            <div className="w-3 h-3 rounded-full bg-red-500"></div>
            <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
            <div className="w-3 h-3 rounded-full bg-green-500"></div>
          </div>
          <div className="flex-1 text-center">
            <span className="text-gray-400 font-mono">jk08y@pro-terminal</span>
          </div>
        </div>
        
        <div className="p-4 h-96 overflow-y-auto font-mono text-sm">
          <div className="text-green-400 mb-4">Welcome to DevFun Terminal v1.0.0 🚀
            Type 'help' for available commands</div>
          
          {commands.map((cmd, index) => (
            <div key={index} className="mb-2">
              {cmd.type === 'input' ? (
                <div className="text-blue-400">
                  <span className="text-pink-400">➜</span> {cmd.text}
                </div>
              ) : (
                <div className="text-green-400 pl-4">{cmd.text}</div>
              )}
            </div>
          ))}
          
          <form onSubmit={handleCommand} className="flex items-center">
            <span className="text-pink-400">➜</span>
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              className="flex-1 ml-2 bg-transparent text-white outline-none"
              autoFocus
            />
          </form>
          <div ref={bottomRef} />
        </div>
      </div>
    </div>
  );
};

export default Terminal;
