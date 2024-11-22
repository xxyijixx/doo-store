import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

import { I18nextProvider } from 'react-i18next'
import i18n from './i18n.ts'



// import  { GlobalStoreProvider } from './GlobalStoreContext.tsx'


// const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

// root.render(
//   <GlobalStoreProvider>
//     <App />
//   </GlobalStoreProvider>
// );



createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <I18nextProvider i18n={i18n}>
  
      <App />
        
    </I18nextProvider>
  </StrictMode>,
)
