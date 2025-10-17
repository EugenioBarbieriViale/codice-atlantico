use actix_web::{get, web, App, HttpServer, Responder, HttpResponse};
use sqlx::PgPool;
use std::env;
use dotenvy::dotenv;
use log::info;

#[get("/healthz")]
async fn healthz(db: web::Data<PgPool>) -> impl Responder {
    if let Err(e) = sqlx::query("SELECT 1").execute(db.get_ref()).await {
        return HttpResponse::InternalServerError().body(format!("DB error: {}", e));
    }
    HttpResponse::Ok().body("ok")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();
    env_logger::init();

    let db_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let port = env::var("CHAT_PORT").unwrap_or_else(|_| "8082".to_string());

    let pool = PgPool::connect(&db_url)
        .await
        .expect("Failed to connect to database");

    info!("Chat service running on port {}", port);

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(pool.clone()))
            .service(healthz)
    })
    .bind(("0.0.0.0", port.parse().unwrap()))?
    .run()
    .await
}
