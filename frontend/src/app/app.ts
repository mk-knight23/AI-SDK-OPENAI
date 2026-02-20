import { Component, signal, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { HttpClientModule } from '@angular/common/http';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, HttpClientModule],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App implements OnInit {
  protected readonly title = signal('MarketPulse');

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.checkApiStatus();
  }

  private checkApiStatus() {
    const apiUrl = '/api/health';
    this.http.get(apiUrl).subscribe({
      next: (response: any) => {
        const statusEl = document.getElementById('api-status');
        if (statusEl) {
          statusEl.textContent = `${response.status} (v${response.version})`;
          statusEl.style.color = '#4ade80';
        }
      },
      error: () => {
        const statusEl = document.getElementById('api-status');
        if (statusEl) {
          statusEl.textContent = 'Offline (development mode)';
          statusEl.style.color = '#fbbf24';
        }
      },
    });
  }
}
