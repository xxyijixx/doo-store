import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'



// import  { GlobalStoreProvider } from './GlobalStoreContext.tsx'


// const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

// root.render(
//   <GlobalStoreProvider>
//     <App />
//   </GlobalStoreProvider>
// );



createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
