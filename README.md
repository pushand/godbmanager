<p># godbmanager</p>
<p>MySql Manager to handle DB transaction for go project</p>
<h2 style="text-align: left;">Start/Stop Service</h2>
<ol>
<li>Start Sql manager by calling StartMySqlService(/* pass MySqlConfig Struct Object that holds sql connection details*/)</li>
<li>Use defer StopMySqlService()</li>
</ol>
<h2>Handling Database Query Service</h2>
<p>&nbsp;</p>
<ol>
<li>Call GetSqlHandler()</li>
<li>The above function will help you perform following
<ul>
<li>Insert - Documentation pending</li>
<li>Update -&nbsp;Documentation pending</li>
<li>QueryRow
<ul>
<li><code>sqlHandler := godbmanager.GetSqlHandler()</code></li>
<li><code>sqlHandler.QueryRow("SELECT name FROM user WHERE id = ?", 1).Scan(&amp;id)</code></li>
</ul>
</li>
<li>AddTransaction - This is when you want to perform bulk transactions
<ul>
<li><code>sqlHandler := godbmanager.GetSqlHandler()</code></li>
<li><code>transactionIdOne := sqlHandler.AddTransaction("INSERT INTO user(name) VALUES (?);", "Gonerd")</code></li>
<li><code>transactionIdTwo := sqlHandler.AddTransaction("INSERT INTO user_detail(id, name) VALUES (?,?);", transactionIdOne, "Gonerd")</code></li>
<li>Above, second transaction is dependent on first transaction id hence we pass transaction id of first transaction</li>
</ul>
</li>
<li>PerformTansactions
<ul>
<li>Once we have added transactions by above AddTransactions, We need to call PerformTransactions() to commit</li>
</ul>
</li>
</ul>
</li>
<li>Thats it</li>
</ol>
