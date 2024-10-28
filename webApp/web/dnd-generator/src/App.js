import './components/Button/Button.js'
import './App.css';
import Button from './components/Button/Button.js';
import Header from './components/Header/Header.js';
import Footer from './components/Footer/Footer.js';

function App() {
  const handleButtonClick = (value) => {
    alert(`Button clicked! Value: ${value}`);
};

return (
    <div>
        <Header />
        <h1>Welcome to My App</h1>
        <Button 
            text="Click Me" 
            value="Button Value 1" 
            onClick={handleButtonClick} 
        />
        <Button 
            text="Another Button" 
            value="Button Value 2" 
            onClick={handleButtonClick} 
        />
        <Footer />

    </div>
);
}

export default App;
